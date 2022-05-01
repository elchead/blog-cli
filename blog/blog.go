package blog

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/afero"
)
const obsidianVault = "/Users/adria/Library/Mobile Documents/iCloud~md~obsidian/Documents/Second_brain"


func GetFilepath(articleTitle,folderPath string) string {
	return path.Join(folderPath,articleTitle+".md")
}

type FsSymLinker interface {
	Symlink(target, symlink string) error
	MkdirAll(path string, perm os.FileMode) error 
}

type Fs interface {
	FsSymLinker
	Create(name string) (afero.File, error)
	Open(name string) (afero.File, error)
}

type Metadata struct {
	Title string
	Date  string
	Categories []string
}

func (m Metadata) String() string {
	return fmt.Sprintf(`---
title: %s
categories: %v
resources:
  - name: "featured-image"
    src: "cover.jpg"
date: %s
---`,m.Title,m.Categories,m.Date)
}



type BlogWriter struct {
	RepoPath string
	FS Fs	
}

func (b *BlogWriter) Write(post Post) error {
	writingFilePath := post.Path()
	file,err := b.FS.Create(writingFilePath)
	if err != nil {
		return errors.Wrapf(err,"could not create post file %s",writingFilePath)
	}
	log.Printf("Created post file: %s", writingFilePath)
	post.Write(file)
	return nil
}

func (b *BlogWriter) AddMedia(post Post,media io.Reader,filename string) error {
	postDir := path.Dir(b.getRepoPostFilePath(post))
	mediaPath := path.Join(postDir,filename)
	file,err := b.FS.Create(mediaPath)	
	if err != nil {
		return errors.Wrapf(err,"could not add media. Are you sure you provided the post type?")
	}
	_, err = io.Copy(file,media)
	return err
}

func (b BlogWriter) LinkInRepo(post Post) error {
	targetFile := post.Path()
	_,openErr := b.FS.Open(targetFile)
	if openErr != nil {
		return errors.Wrap(openErr,"Failed to link file to non existing post")
	}
	symlink := b.getRepoPostFilePath(post)

	err := b.mkdir(path.Dir(symlink))
	if err != nil {
		return err
	}
	err = b.mkdir(path.Dir(targetFile))
	if err != nil {
		return err
	}

	err = b.FS.Symlink(targetFile,symlink)
	if err != nil {
		log.Println(err)	
	}
	return nil
}

func (b BlogWriter) mkdir(path string) error {
	err := b.FS.MkdirAll(path,0777)
	if err != nil {
		return fmt.Errorf("could not create directory: %w", err)
	}
	log.Printf("Created directory: %s", path)	
	return nil
}

func (b BlogWriter) getRepoPostFilePath(post Post) string {
	return ConstructRepoPostFilePath(b.RepoPath,post.RepoFolder(),post.Title())
}

func constructDirNameFromTitle(title string) string {
	lowerCase := strings.ToLower(title)
	cutAfterDash := strings.Split(lowerCase," - ")[0]
	noSpaces := strings.Replace(cutAfterDash, " ","-",-1)
	return noSpaces
}

// get the (sub-)link to the blog post
func ConstructPostLink(post Post) string {
	return path.Join(post.RepoFolder(),constructDirNameFromTitle(post.Title()))
}

func ConstructRepoPostFilePath(repoPath ,postType, postTitle string) string {
	return path.Join(repoPath,"content",postType,constructDirNameFromTitle(postTitle),"index.en.md")
}







