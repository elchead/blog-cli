package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/elchead/blog-cli/blog"
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

func main() {

	app := &cli.App{
		Name:  "Blog-CLI",
		Authors: []*cli.Author{{Name:"Adrian Stobbe",Email:"stobbe.adrian@gmail.com"}},
		Usage: "quickly generate blog posts inside Obsidian and publish on Github with Hugo",
		Commands: []*cli.Command{
			{
				Name: "post",
				Usage: "create new post",
				Action: func(c *cli.Context) error {
					title := c.Args().Get(0)
					fmt.Printf("Create new post %s\n",title)
					fmt.Print("Enter category: ")
					reader := bufio.NewReader(os.Stdin)
					category, err := reader.ReadString('\n')
					if err != nil {
						fmt.Println("An error occured while reading input. Please try again", err)
						return nil
					}
					category = strings.TrimSuffix(category, "\n")
					meta := blog.Metadata{Title: title, Categories : []string{category}, Date: time.Now().Format("2006-01-02")}
					writingPath := blog.GetFilepath(meta.Title,writingDir)
					file,err := os.Create(writingPath)
					if err != nil {
						log.Fatal(err)
					}
					b := blog.Blog{RepoPath:repoDir}
					b.WritePost(meta,file)
					err = b.CreatePostInRepo(fs,meta,writingPath)
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
