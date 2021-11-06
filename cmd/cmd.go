package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/elchead/blog-cli/blog"
	"github.com/skratchdot/open-golang/open"
	"github.com/urfave/cli/v2"
)

type Fs struct {}

func (f Fs) Symlink(target,link string) error {
	return os.Symlink(target,link)
}
func (f Fs) MkdirAll(path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
}

const writingDir = "/Users/adria/Google Drive/Obsidian/Second_brain/Blog"
const repoDir = "/Users/adria/Programming/elchead.github.io"
var fs = Fs{}

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

func createWriterFile(title, writingPath string) *os.File {
	file,err := os.Create(writingPath)
	if err != nil {
		log.Fatal(err)
	}
	return file	
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
				Action: func(c *cli.Context) error {
					meta := readMetadata(c.Args().Get(0))
					b := blog.Blog{RepoPath:repoDir}
					writingFilePath := blog.GetFilepath(meta.Title,writingDir)
					b.WritePost(meta,createWriterFile(meta.Title,writingFilePath))
					err := b.CreatePostInRepo(fs,meta.Title,writingFilePath)
					if err != nil {
						log.Fatal(err)
					}
					OpenObsidianFile(filepath.Base(writingFilePath))
					return nil
				},
			},
			{
				Name: "draft",
				Usage: "create new post without reference in repo",
				Action: func(c *cli.Context) error {
					meta := readMetadata(c.Args().Get(0))
					b := blog.Blog{RepoPath:repoDir}
					writingFilePath := blog.GetFilepath(meta.Title,writingDir)
					b.WritePost(meta,createWriterFile(meta.Title,writingFilePath))
					OpenObsidianFile(filepath.Base(writingFilePath))
					return nil
				},
			},
			{
				Name: "publish",
				Usage: "use existing obsidian file to create reference in repo",
				Action: func(c *cli.Context) error {
					title := c.Args().Get(0)
					b := blog.Blog{RepoPath:repoDir}
					writingFilePath := blog.GetFilepath(title,writingDir)
					err := b.CreatePostInRepo(fs,title,writingFilePath)
					if err != nil {
						log.Fatal(err)
					}
					return nil
				},
			},
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func OpenObsidianFile(filename string) {
	err := open.Run(fmt.Sprintf("obsidian://open?file=%s",filename))
	if err != nil {
		log.Printf("Error opening obsidian: %v", err)
	}
}
