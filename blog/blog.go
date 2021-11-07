package blog

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"strings"

	"github.com/spf13/afero"
)

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
date: %s
---`,m.Title,m.Categories,m.Date)
}

type Blog struct {
	RepoPath string
	WritingDir string
	FS Fs	
}

type Post interface {
	Title() string
	Write(file io.Writer) // error
}

func CreatePost(meta Metadata,file afero.File, writingFilePath string) Post {
	if meta.Categories[0] == "Book-notes" {
		return Book{TemplateFile: file, Meta: meta}
	} else {
		return Article{Meta: meta,File: file, Path: writingFilePath}
	}
}

func (b *Blog) DraftPost(meta Metadata) (Article,error) {
	writingFilePath := GetFilepath(meta.Title,b.WritingDir)
	file,err := b.FS.Create(writingFilePath)
	if err != nil {
		return Article{},err
	}
	post := Article{Meta: meta,File: file, Path: writingFilePath} //CreatePost(meta,file, writingFilePath)
	post.Write(file)
	return post,nil
}

func (b Blog) LinkInRepo(article Post) error {
	title := article.Title()
	return b.LinkInRepoFromTitle(title)
	
}

func (b Blog) mkdir(path string) error {
	err := b.FS.MkdirAll(path,0777)
	if err != nil {
		return fmt.Errorf("could not create directory: %w", err)
	}
	log.Printf("Created directory: %s", path)	
	return nil
}

func (b Blog) LinkInRepoFromTitle(title string) error {
	targetFile := GetFilepath(title,b.WritingDir) 
	symlink := b.getSimpleRepoPostFilePath(title)

	err := b.mkdir(path.Dir(symlink))
	if err != nil {
		return err
	}
	err = b.mkdir(path.Dir(targetFile))
	if err != nil {
		return err
	}

	return b.FS.Symlink(targetFile,symlink)
}

func (b Blog) getSimpleRepoPostFilePath(title string) string {
	return constructRepoPostFilePath(b.RepoPath,title)
}

type Article struct {
	Meta Metadata
	File io.Writer
	Path string
}

func (a Article) Title() string {
	return a.Meta.Title
}

func (b Article) Write(file io.Writer) {
	io.WriteString(file,b.Meta.String())
}

func constructDirNameFromTitle(title string) string {
	lowerCase := strings.ToLower(title)
	cutAfterDash := strings.Split(lowerCase," - ")[0]
	noSpaces := strings.Replace(cutAfterDash, " ","-",-1)
	return noSpaces
}

func constructRepoPostFilePath(repoPath ,dirName string) string {
	return path.Join(repoPath,"content","posts",constructDirNameFromTitle(dirName),"index.en.md")
}







