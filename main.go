package main

import (
	"log"
	"os"

	"github.com/elchead/blog-cli/blog"
)

func main() {
	b := blog.Blog{}
	meta := blog.Metadata{Title: "title", Categories : []string{"Thoughts"}, Date: "2021-11-04"}
	file,err := os.Create(blog.GetFilepath(meta.Title,"/Users/adria/Google Drive/Obsidian/Second_brain/Blog"))
	if err != nil {
		log.Fatal(err)
	}
	b.WritePost(meta,file)
}
