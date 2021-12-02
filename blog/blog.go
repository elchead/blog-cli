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


type Post interface {
	Title() string
	Write(file io.Writer)
	RepoFolder() string
	Path() string
}
type BlogWriter struct {
	RepoPath string
	WritingDir string

	BookDir string
	BookTemplate io.Reader

	FS Fs	
}

func (b *BlogWriter) DraftArticle(meta Metadata) (Article,error) {
	writingFilePath := GetFilepath(meta.Title,b.WritingDir)
	file,err := b.FS.Create(writingFilePath)
	if err != nil {
		return Article{},err
	}
	post := Article{Meta: meta,File: file, path: writingFilePath}
	post.Write(file)
	return post,nil
}

func (b *BlogWriter) DraftBook(meta Metadata) (Book,error) {
	if b.BookDir == "" || b.BookTemplate == nil {
		log.Fatal("Define book parameters before drafting a book")
	}
	writingFilePath := GetFilepath(meta.Title,b.BookDir)
	file,err := b.FS.Create(writingFilePath)
	if err != nil {
		return Book{},errors.Wrapf(err,"could not create book file %s",writingFilePath)
	}
	log.Printf("Created book file: %s", writingFilePath)
	post := Book{b.BookTemplate,meta,writingFilePath}
	post.Write(file)
	return post,nil
}

func (b *BlogWriter) AddMedia(post Post,media io.Reader,filename string) error {
	postDir := path.Dir(b.getRepoPostFilePath(post))
	mediaPath := path.Join(postDir,filename)
	file,err := b.FS.Create(mediaPath)	
	if err != nil {
		errors.Wrapf(err,"could not add media to %s",mediaPath)
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

func ConstructRepoPostFilePath(repoPath ,postType, dirName string) string {
	return path.Join(repoPath,"content",postType,constructDirNameFromTitle(dirName),"index.en.md")
}







