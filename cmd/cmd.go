package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/elchead/blog-cli/blog"
	"github.com/elchead/blog-cli/fs"
	"github.com/elchead/blog-cli/git"
	"github.com/urfave/cli/v2"
)

const writingDir = "/Users/adria/Google Drive/Obsidian/Second_brain/Blog"
const repoDir = "/Users/adria/Programming/elchead.github.io"
const mediaDir = "/Users/adria/Downloads"
const bookDir = "/Users/adria/Google Drive/Obsidian/Second_brain/Books"
const bookTemplatePath = "/Users/adria/Google Drive/Obsidian/Second_brain/Templates/book.md"
var filesystem = fs.Filesystem{}

var blogWriter = createBlog()

var bookFlag = &cli.BoolFlag{
	Name: "book",
	Aliases: []string{"B"},
	Value: false,
	Usage: "set if post is book-note",
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
					if(c.Bool("book")){
						AskToPublishToReadwise(os.Stdin,post,PushToReadwise)
					}
					return blogPusher.Push(post)
				},
			},
			{
				Name: "readwise",
				Usage: "Push book notes to readwise",
				ArgsUsage: "provide topic of post. Assumes that the post is linked in the repository",
				Action: func(c *cli.Context) error {
					post := newPost(c.Args().Get(0),true)
					AskToPublishToReadwise(os.Stdin,post,PushToReadwise)
					return nil
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
					media, err := filesystem.Open(path.Join(mediaDir,mediaFilename))
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

func AskToPublishToReadwise(read io.Reader, post blog.Post,push func(path string)){
	isYes := getInput(read,"Do you want to publish the book note? (y!/n!)")
	if(isYes){
		postPath := filepath.Dir(blog.ConstructRepoPostFilePath(repoDir,post.RepoFolder(),post.Title()))
		push(postPath)
	} else {
		fmt.Println("Not publishing to readwise")
	}
}

func PublishIfInputYes(post blog.Post) {
	if okToPublish(os.Stdin) {
		blogPusher := git.NewBlogPush(blogWriter.RepoPath)
		blogPusher.Push(post)
	}
}

func okToPublish(read io.Reader) bool {
	return getInput(read,"Publish post (y!/n!): ")	
}

