package main

import (
	"log"
	"os"

	"github.com/elchead/blog-cli/blog"
)

func main() {
	b := blog.Blog{}
	meta := blog.Metadata{Title: "title", Categories : []string{"Thoughts"}, Date: "2021-11-04"}
	fileHelper := blog.File{Title: "title", Path:"."}
	file,err := os.Create(fileHelper.Filepath())
	if err != nil {
		log.Fatal(err)
	}
	b.WritePost(meta,file)
}
