package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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
func draftPost(title string,isBook bool) blog.Post {
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
func createPost(title string,isBook bool) blog.Post {
	var post blog.Post
	meta := blog.Metadata{Title:title}
	if isBook {
		post = blog.Book{Meta:meta,Path_:blog.GetFilepath(title,bookDir)}
	} else {
		post = blog.Article{Meta:meta,Path_: blog.GetFilepath(title,writingDir)}
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
				Flags: []cli.Flag{
					bookFlag,
				},
				Action: func(c *cli.Context) error {
					post := draftPost(c.Args().Get(0), c.Bool("book"))
					blogger.LinkInRepo(post)
					OpenObsidianFile(filepath.Base(post.Path()))	
					return nil
				},
			},
			{
				Name: "draft",
				Usage: "create new post without reference in repo",
				Flags: []cli.Flag{
					bookFlag,
				},
				Action: func(c *cli.Context) error {
					post := draftPost(c.Args().Get(0), c.Bool("book"))
					OpenObsidianFile(filepath.Base(post.Path()))	
					return nil
				},
			},
			{
				Name: "publish",
				Usage: "use existing obsidian article to create reference in repo. Then open preview",
				Flags: []cli.Flag{
					bookFlag,
				},
				Action: func(c *cli.Context) error {
					post := createPost(c.Args().Get(0),c.Bool("book"))
					err := blogger.LinkInRepo(post)
					if err != nil {
						log.Fatal(err)
					}
					ok := OpenBrowser()
					RenderBlog(ok)
					return nil
				},
			},
			{
				Name: "preview",
				Usage: "render blog and open",
				Action: func(c *cli.Context) error {
					ok := OpenBrowser()
					RenderBlog(ok)
					return nil
				},
			},
			{
				Name: "push",
				Usage: "render blog and open",
				Action: func(c *cli.Context) error {
					post := createPost(c.Args().Get(0),c.Bool("book"))
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

func RenderBlog(b BrowserOk) {
	cmd := exec.Command("hugo","serve")
	cmd.Dir = repoDir
	err := cmd.Run()
	if err != nil {
		log.Fatal("Could not serve hugo: ",err)
	}
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
