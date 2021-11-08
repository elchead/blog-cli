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

	BookDir string // todo??
	BookTemplate io.Reader

	FS Fs	
}

type Post interface {
	Title() string
	Write(file io.Writer) // error
}

func  (b Blog) GetFilepathFromMeta(meta Metadata) string{
	if meta.Categories[0] == "Book-notes" {
		return GetFilepath(meta.Title,b.BookDir)
	
	} else {
		return GetFilepath(meta.Title,b.WritingDir)
	}	
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

func (b *Blog) DraftBook(meta Metadata) (Book,error) {
	if b.BookDir == "" || b.BookTemplate == nil {
		log.Fatal("Define book parameters before drafting a book")
	}
	writingFilePath := GetFilepath(meta.Title,b.BookDir)
	file,err := b.FS.Create(writingFilePath)
	if err != nil {
		return Book{},errors.Wrapf(err,"could not create book file %s",writingFilePath)
	}
	log.Printf("Created book file: %s", writingFilePath)
	post := Book{b.BookTemplate,meta}
	post.Write(file)
	return post,nil
}

func (b Blog) LinkInRepo(article Post) error {
	title := article.Title()
	targetFile := GetFilepath(title,b.WritingDir) 
	symlink := b.getSimpleRepoPostFilePath(article)

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

func (b Blog) mkdir(path string) error {
	err := b.FS.MkdirAll(path,0777)
	if err != nil {
		return fmt.Errorf("could not create directory: %w", err)
	}
	log.Printf("Created directory: %s", path)	
	return nil
}

func (b Blog) getSimpleRepoPostFilePath(post Post) string {
	var postType string
	switch post.(type) {
	case Book:
		postType = "books"
	case Article:
		postType = "posts"
	}
	return constructRepoPostFilePath(b.RepoPath,postType,post.Title())
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

func constructRepoPostFilePath(repoPath ,postType, dirName string) string {
	return path.Join(repoPath,"content",postType,constructDirNameFromTitle(dirName),"index.en.md")
}







