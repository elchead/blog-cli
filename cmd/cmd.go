package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/elchead/blog-cli/blog"
	"github.com/elchead/blog-cli/git"
	"github.com/skratchdot/open-golang/open"
	"github.com/spf13/afero"
	"github.com/urfave/cli/v2"
)

const writingDir = "/Users/adria/Google Drive/Obsidian/Second_brain/Blog"
const repoDir = "/Users/adria/Programming/elchead.github.io"
const mediaDir = "/Users/adria/Downloads"
const bookDir = "/Users/adria/Google Drive/Obsidian/Second_brain/Books"
const bookTemplatePath = "/Users/adria/Google Drive/Obsidian/Second_brain/Templates/book.md"
var fs = Filesystem{}

var blogWriter = createBlog()

var bookFlag = &cli.BoolFlag{
	Name: "book",
	Aliases: []string{"B"},
	Value: false,
	Usage: "set if post is book-note",
      }


// instantiate and draft Post
func createAndWritePost(title string,isBook bool) blog.Post {
	var post blog.Post
	var err error
	if !isBook {
		meta := readMetadata(title)
		post,err = blogWriter.DraftArticle(meta)
	} else {
		bookmeta := blog.Metadata{Title: title}
		post, err = blogWriter.DraftBook(bookmeta)
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
				Usage: "create new post with reference in repo",
				ArgsUsage: "provide post topic (used for folder and file naming)",
				Flags: []cli.Flag{
					bookFlag,
				},
				Action: func(c *cli.Context) error {
					if c.Args().Len() == 0 {
						return fmt.Errorf("no title specified")
					}
					post := createAndWritePost(c.Args().First(), c.Bool("book"))
					blogWriter.LinkInRepo(post)
					OpenObsidianFile(filepath.Base(post.Path()))	
					return nil
				},
			},
			{
				Name: "draft",
				Usage: "create new post without reference in repo",
				ArgsUsage: "provide post topic (used for folder and file naming)",
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
				Usage:"use existing Obsidian article to create linkage in repo. Then locally render blog (`hugo serve`) and open preview in Browser. Finally, it asks if you want to publish the post.",
				ArgsUsage: "provide title of existing Obsidian file",
				Flags: []cli.Flag{
					bookFlag,
				},
				Action: func(c *cli.Context) error {
					if c.Args().Len() == 0 {
						return fmt.Errorf("no title specified")
					}
					post := newPost(c.Args().Get(0),c.Bool("book"))
					err := blogWriter.LinkInRepo(post)
					if err != nil {
						log.Fatal(err)
					}
					link :=  blog.ConstructPostLink(post)
					cmd := OpenPostInBrowser(link)
					PublishIfInputYes(post)
					cmd.Wait()
					return nil
				},
			},
			{
				Name: "preview",
				Usage: "render blog and open",
				Action: func(c *cli.Context) error {
					cmd := OpenPostInBrowser("")
					cmd.Wait()
					return nil
				},
			},
			{
				Name: "push",
				Usage: "handles git logic for publishing. It stages existing changes, replaces the symbolic link with a hard link, commits, pulls and pushes.",
				ArgsUsage: "provide topic of post. Assumes that the post is linked in the repository",
				Flags: []cli.Flag{
					bookFlag,
				},
				Action: func(c *cli.Context) error {
					post := newPost(c.Args().Get(0),c.Bool("book"))
					blogPusher := git.NewBlogPush(blogWriter.RepoPath)
					return blogPusher.Push(post)
				},
			},
			{
				Name: "media",
				Usage: "add media to post. Copies file to git repository",
				ArgsUsage: "provide topic of post. Assumes that the post is linked in the repository. Second argument is media filename inside media directory",
				Flags: []cli.Flag{
					bookFlag,
				},
				Action: func(c *cli.Context) error {
					if c.Args().Len() < 2 {
						return fmt.Errorf("please specify topic and media filename")
					}
					post := newPost(c.Args().Get(0),c.Bool("book"))
					mediaFilename := c.Args().Get(1)
					media, err := fs.Open(path.Join(mediaDir,mediaFilename))
					if err != nil {
						log.Fatalf("Media could not be opened: %v", err)
					}
					return blogWriter.AddMedia(post,media,mediaFilename)
				},
			},
		},

	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func OpenPostInBrowser(link string) *exec.Cmd {
	cmd := StartRenderBlog()
	time.Sleep(1 * time.Second)
	OpenBrowser(link)
	fmt.Println("Press Ctrl+c to stop render process")
	return cmd
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
		blogPusher := git.NewBlogPush(blogWriter.RepoPath)
		blogPusher.Push(post)
	}
}

func OpenBrowser(link string) {
	url := "http://localhost:1313/"+link
	err := open.Run(url)
	if err != nil {
		fmt.Println("Could not open browser: ", err)
	}
	fmt.Println("Open ",url)
}

func StartRenderBlog() *exec.Cmd {
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

func createBlog() blog.BlogWriter {
	templateFile, err := os.Open(bookTemplatePath)
	if err != nil {
		log.Fatal(err)
	}
	return blog.BlogWriter{RepoPath:repoDir,WritingDir: writingDir,FS:fs,BookDir:bookDir,BookTemplate:templateFile}
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
		cmd := StartRenderBlog()		
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
