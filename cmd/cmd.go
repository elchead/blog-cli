package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/elchead/blog-cli/blog"
	"github.com/skratchdot/open-golang/open"
	"github.com/spf13/afero"
	"github.com/urfave/cli/v2"
)

const writingDir = "/Users/adria/Google Drive/Obsidian/Second_brain/Blog"
const repoDir = "/Users/adria/Programming/elchead.github.io"
const bookDir = "/Users/adria/Google Drive/Obsidian/Second_brain/Books"
const bookTemplatePath = "/Users/adria/Google Drive/Obsidian/Second_brain/Templates/book.md"
var fs = Filesystem{}

var blogger = createBlog()

var bookFlag = &cli.BoolFlag{
	Name: "book",
	Aliases: []string{"B"},
	Value: false,
	Usage: "set to create book template",
      }


// instantiate and draft Post
func createAndWritePost(title string,isBook bool) blog.Post {
	var post blog.Post
	var err error
	if !isBook {
		meta := readMetadata(title)
		post,err = blogger.DraftArticle(meta)
	} else {
		bookmeta := blog.Metadata{Title: title}
		post, err = blogger.DraftBook(bookmeta)
		fmt.Printf("Created new book note %s\n",title)
	}
	if err != nil {
		log.Fatal(err)
	}
	return post
}

// only instantiate Post
func newPost(title string,isBook bool) blog.Post {
	var post blog.Post
	meta := blog.Metadata{Title:title}
	if isBook {
		post = blog.NewBook(meta, blog.GetFilepath(title,bookDir))
	} else {
		post = blog.NewArticle(meta,blog.GetFilepath(title,writingDir))
	}
	return post
}



func main() {

	app := &cli.App{
		Name:  "Blog-CLI",
		Authors: []*cli.Author{{Name:"Adrian Stobbe",Email:"stobbe.adrian@gmail.com"}},
		Usage: "quickly generate blog posts inside Obsidian and publish on Github with Hugo",
		Commands: []*cli.Command{
			{
				Name: "post",
				Usage: "provide post title",
				Description: "create new post with reference in repo",
				Flags: []cli.Flag{
					bookFlag,
				},
				Action: func(c *cli.Context) error {
					if c.Args().Len() == 0 {
						return fmt.Errorf("no title specified")
					}
					post := createAndWritePost(c.Args().First(), c.Bool("book"))
					blogger.LinkInRepo(post)
					OpenObsidianFile(filepath.Base(post.Path()))	
					return nil
				},
			},
			{
				Name: "draft",
				Description: "create new post without reference in repo",
				Usage: "provide post title",
				Flags: []cli.Flag{
					bookFlag,
				},
				Action: func(c *cli.Context) error {
					if c.Args().Len() == 0 {
						return fmt.Errorf("no title specified")
					}
					post := createAndWritePost(c.Args().Get(0), c.Bool("book"))
					OpenObsidianFile(filepath.Base(post.Path()))	
					return nil
				},
			},
			{
				Name: "preview-post",
				Description:"use existing Obsidian article to create linkage in repo. Then locally render blog (`hugo serve`) and open preview in Browser. Finally, it asks if you want to publish the post.",
				Usage: "provide title of existing Obsidian file",
				Flags: []cli.Flag{
					bookFlag,
				},
				Action: func(c *cli.Context) error {
					if c.Args().Len() == 0 {
						return fmt.Errorf("no title specified")
					}
					post := newPost(c.Args().Get(0),c.Bool("book"))
					err := blogger.LinkInRepo(post)
					if err != nil {
						log.Fatal(err)
					}
					ok := OpenBrowser()
					cmd := StartRenderBlog(ok)
					PublishIfInputYes(post)
					fmt.Println("Press Ctrl+c to stop render process")
					cmd.Wait()
					return nil
				},
			},
			{
				Name: "preview",
				Description: "render blog and open",
				Action: func(c *cli.Context) error {
					ok := OpenBrowser()
					cmd := StartRenderBlog(ok)
					fmt.Println("Press Ctrl+c to stop render process")
					cmd.Wait()
					return nil
				},
			},
			{
				Name: "push",
				Description: "handles git logic for publishing. It stages existing changes, replaces the symbolic link with a hard link, commits, pulls and pushes.",
				Usage: "provide title of post. Assumes that the post is linked in the repository",
				Action: func(c *cli.Context) error {
					post := newPost(c.Args().Get(0),c.Bool("book"))
					return blogger.Push(post)
				},
			},
		},

	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func okToPublish(read io.Reader) bool {
	fmt.Print("Publish post (y!/n!): ")
	reader := bufio.NewReader(read)
	answer, _ := reader.ReadString('\n')
	answer = strings.TrimSuffix(answer, "!\n")
	desiredInput := "y" 
	if answer != desiredInput {
		fmt.Println("Answer was:", answer)
	}
	return answer == desiredInput		
} 

func PublishIfInputYes(post blog.Post) {
	if okToPublish(os.Stdin) {
		blogger.Push(post)
	}
}

type BrowserOk struct {}
// need to open browser before rendering, since render remains active
func OpenBrowser() BrowserOk {
	url := "http://localhost:1313/"
	err := open.Run(url)
	if err != nil {
		fmt.Println("Could not open browser: ", err)
	}
	fmt.Println("Open ",url)
	return BrowserOk{}
}

func StartRenderBlog(b BrowserOk) *exec.Cmd {
	cmd := exec.Command("hugo","serve")//"--disableFastRender"
	cmd.Dir = repoDir
	err := cmd.Start()
	output, _ := cmd.CombinedOutput()
	if std:=string(output); std!= "" { 
		fmt.Println(std) 
	}
	if err != nil {
		log.Fatal("Could not serve hugo: ",err)
	}
	return cmd
}

func OpenObsidianFile(filename string) {
	err := open.Run(fmt.Sprintf("obsidian://open?file=%s",filename))
	if err != nil {
		log.Printf("Error opening obsidian: %v", err)
	}
}

func createBlog() blog.Blog {
	templateFile, err := os.Open(bookTemplatePath)
	if err != nil {
		log.Fatal(err)
	}
	return blog.Blog{RepoPath:repoDir,WritingDir: writingDir,FS:fs,BookDir:bookDir,BookTemplate:templateFile}
}

func readMetadata(title string) blog.Metadata {
					fmt.Printf("Create new post %s\n",title)
					fmt.Print("Enter category: ")
					reader := bufio.NewReader(os.Stdin)
					category, err := reader.ReadString('\n')
					if err != nil {
						log.Fatal("An error occured while reading input. Please try again", err)
					}
					category = strings.TrimSuffix(category, "\n")
					return blog.Metadata{Title: title, Categories : []string{category}, Date: time.Now().Format("2006-01-02")}
}

type Filesystem struct {}

func (f Filesystem) Symlink(target,link string) error {
	return os.Symlink(target,link)
}
func (f Filesystem) MkdirAll(path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
}
func (f Filesystem) Create(path string) (afero.File,error) {
	return os.Create(path)
}

func (f Filesystem) Open(path string) (afero.File,error) {
	return os.Open(path)
}

func startGoRoutine(exitChan chan os.Signal, done chan bool) {
	go func() {
		cmd := StartRenderBlog(BrowserOk{})		
		for s := range exitChan {
			switch s {
			case syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT,os.Interrupt:
				done <- true
				return
			}
		}
		cmd.Wait()
	}()
}
