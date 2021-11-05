package main

import (
	"log"
	"os"

	"github.com/elchead/blog-cli/blog"
)

type Fs struct {}

func (f Fs) Symlink(target,link string) error {
	return os.Symlink(target,link)
}
func (f Fs) MkdirAll(path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
}

func main() {
	b := blog.Blog{RepoPath:"/Users/adria/Programming/elchead.github.io"}
	meta := blog.Metadata{Title: "title", Categories : []string{"Thoughts"}, Date: "2021-11-04"}
	originalFpath := blog.GetFilepath(meta.Title,"/Users/adria/Google Drive/Obsidian/Second_brain/Blog")
	file,err := os.Create(originalFpath)
	if err != nil {
		log.Fatal(err)
	}
	b.WritePost(meta,file)
	fs := Fs{}
	err  = b.CreatePost(fs,meta,originalFpath)
	if err != nil {
		log.Fatal(err)
	}
	
}
