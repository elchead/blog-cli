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
	"github.com/elchead/blog-cli/fs"
	"github.com/elchead/blog-cli/git"
	"github.com/skratchdot/open-golang/open"
	"github.com/urfave/cli/v2"
)

const repoDir = "/Users/adria/Programming/elchead.github.io"
const mediaDir = "/Users/adria/Downloads"
const obsidianVault = "/Users/adria/Library/Mobile Documents/iCloud~md~obsidian/Documents/Second_brain"
const writingDir = obsidianVault +"/Blog"
const bookDir = obsidianVault +"/Books"
const bookTemplatePath = obsidianVault +"/Templates/book.md"

var filesystem = fs.Filesystem{}

var blogWriter = newBlog()
var postFactory = newPostFactory()

var bookFlag = &cli.BoolFlag{
	Name: "book",
	Aliases: []string{"B"},
	Value: false,
	Usage: "set if post is book-note",
      }

var letterFlag = &cli.BoolFlag{
	Name: "letter",
	Aliases: []string{"L"},
	Value: false,
	Usage: "set if post is letter",
}

var flags = []cli.Flag{
	bookFlag,letterFlag,
}

func newPostFactory() blog.PostFactory {
	templateFile, err := os.Open(bookTemplatePath)
	if err != nil {
		log.Fatal(err)
	}
	return blog.PostFactory{BookTemplate:templateFile,BaseDir: obsidianVault}
}

func newMetadata(title string,isBook,isLetter bool) blog.Metadata {
	return newMetadataFrom(title,isLetter,isBook,os.Stdin)
}

func newMetadataFrom(title string,isLetter,isBook bool,input io.Reader) blog.Metadata {
	if isBook {
		return blog.Metadata{Title: title, Categories : []string{"Book-notes"}}
	} else if isLetter {
		return blog.Metadata{Title: title, Categories : []string{"Letters"}, Date: time.Now().Format("2006-01-02")}
	} else {
		return readMetadataFrom(title,input)
	}
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
				ArgsUsage: "provide post name (used for folder and file naming)",
				Flags: flags,
				Action: func(c *cli.Context) error {
					if c.Args().Len() == 0 {
						return fmt.Errorf("no title specified")
					}
					meta := newMetadata(c.Args().First(), c.Bool("book"),c.Bool("letter"))
					post,err := postFactory.NewPost(meta)
					if err != nil {
						log.Fatal(err)
					}
					err = blogWriter.Write(post)
					if err != nil {
						log.Fatal(err)
					}
					err = blogWriter.LinkInRepo(post)
					if err != nil {
						log.Fatal(err)
					}
					OpenObsidianFile(post.Path())	
					return nil
				},
			},
			{
				Name: "draft",
				Usage: "create new post without reference in repo",
				ArgsUsage: "provide post topic (used for folder and file naming)",
				Flags: flags,
				Action: func(c *cli.Context) error {
					if c.Args().Len() == 0 {
						return fmt.Errorf("no title specified")
					}
					meta := newMetadata(c.Args().First(), c.Bool("book"),c.Bool("letter"))
					post,err := postFactory.NewPost(meta)
					if err != nil {
						log.Fatal(err)
					}
					err = blogWriter.Write(post)
					if err != nil {
						log.Fatal(err)
					}
					OpenObsidianFile(post.Path())	
					return nil
				},
			},
			{
				Name: "preview-post",
				Usage:"use existing Obsidian article to create linkage in repo. Then locally render blog (`hugo serve`) and open preview in Browser. Finally, it asks if you want to publish the post.",
				ArgsUsage: "provide title of existing Obsidian file",
				Flags: flags,
				Action: func(c *cli.Context) error {
					if c.Args().Len() == 0 {
						return fmt.Errorf("no title specified")
					}
					meta := newMetadata(c.Args().First(), c.Bool("book"),c.Bool("letter"))
					post,err := postFactory.NewPost(meta)
					if err != nil {
						log.Fatal(err)
					}
					err = blogWriter.LinkInRepo(post)
					if err != nil {
						log.Fatal(err)
					}
					link := blog.ConstructPostLink(post)
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
				Flags: flags,
				Action: func(c *cli.Context) error {
					meta := newMetadata(c.Args().First(), c.Bool("book"),c.Bool("letter"))
					post,err := postFactory.NewPost(meta)
					if err != nil {
						log.Fatal(err)
					}
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
				Flags: []cli.Flag{
					bookFlag,
				},
				Action: func(c *cli.Context) error {
					meta := newMetadata(c.Args().First(), c.Bool("book"),c.Bool("letter"))
					post,err := postFactory.NewPost(meta)
					if err != nil {
						log.Fatal(err)
					}
					AskToPublishToReadwise(os.Stdin,post,PushToReadwise)
					return nil
				},
			},
			{
				Name: "media",
				Usage: "add media to post. Copies file to git repository",
				ArgsUsage: "provide topic of post. Assumes that the post is linked in the repository. Second argument is media filename inside media directory",
				Flags: flags,
				Action: func(c *cli.Context) error {
					if c.Args().Len() < 2 {
						return fmt.Errorf("please specify topic and media filename")
					}
					meta := newMetadata(c.Args().First(), c.Bool("book"),c.Bool("letter"))
					post,err := postFactory.NewPost(meta)
					if err != nil {
						log.Fatal(err)
					}
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
	isYes := getInput(read,"Do you want to publish the book notes on Readwise? (y!/n!)")
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

func getInput(read io.Reader,inputQuestion string) bool {
	fmt.Print(inputQuestion)
	reader := bufio.NewReader(read)
	answer, _ := reader.ReadString('\n')
	answer = strings.TrimSuffix(answer, "!\n")
	desiredInput := "y" 
	if answer != desiredInput {
		fmt.Println("Answer was:", answer)
	}
	return answer == desiredInput		
}

func OpenPostInBrowser(link string) *exec.Cmd {
	OpenBrowser(link)
	cmd := StartRenderBlog()
	time.Sleep(1 * time.Second)
	fmt.Println("Press Ctrl+c to stop render process")
	return cmd
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
	output, err := cmd.CombinedOutput()
	if std:=string(output); std!= "" { 
		fmt.Println(std) 
	}
	if err != nil {
		log.Fatal("Could not serve hugo: ",err)
	}
	return cmd
}

func OpenObsidianFile(path string) {
	err := open.Run(fmt.Sprintf("obsidian://open?file=%s",GetVaultPath(path)))
	if err != nil {
		log.Printf("Error opening obsidian: %v", err)
	}
}

func newBlog() blog.BlogWriter {
	return blog.BlogWriter{RepoPath:repoDir,FS:filesystem}
}

func readMetadata(title string) blog.Metadata {
	return readMetadataFrom(title,os.Stdin)
}

func readMetadataFrom(title string,input io.Reader) blog.Metadata {
	fmt.Printf("Create new post %s\n",title)
	fmt.Print("Enter category: ")
	reader := bufio.NewReader(input)
	category, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal("An error occured while reading input. Please try again", err)
	}
	category = strings.TrimSuffix(category, "\n")
	return blog.Metadata{Title: title, Categories : []string{category}, Date: time.Now().Format("2006-01-02")}
}

// Requires https://github.com/elchead/readwise-note-extractor with inclusion in PATH
func PushToReadwise(path string){
	cmd := exec.Command("push_readwise.py")
	cmd.Dir = path
	output, err := cmd.CombinedOutput()
	if std:=string(output); std!= "" { 
		fmt.Println(std) 
	}
	if err != nil {
		log.Fatal("Could not push to readwise: ",err)
	}
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

func getFolder(path string) (string) {
	return filepath.Base(filepath.Dir(path))
}

func getFilename(path string) (string) {
	return filepath.Base(path)	
}

func GetVaultPath(path string) (string) {
	return getFolder(path)+"/"+getFilename(path)
}
