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

func (b *Blog) DraftPost(meta Metadata) (Article,error) {
	writingFilePath := GetFilepath(meta.Title,b.WritingDir)
	file,err := b.FS.Create(writingFilePath)
	if err != nil {
		return Article{},err
	}
	article := Article{Meta: meta,File: file}
	article.Write(file) // TODO refactor
	return article,nil
}

func (b Blog) LinkInRepo(article Article) error {
	title := article.Meta.Title
	return b.LinkInRepoFromTitle(title)
	
}

func (b Blog) LinkInRepoFromTitle(title string) error {
	targetFile := GetFilepath(title,b.WritingDir) 
	symlink := b.getSimpleRepoPostFilePath(title)
	err := b.FS.MkdirAll(path.Dir(symlink),0777)
	if err != nil {
		return fmt.Errorf("could not create directory: %w", err)
	}
	err = b.FS.MkdirAll(path.Dir(targetFile),0777)
	if err != nil {
		return fmt.Errorf("could not create directory: %w", err)
	}
	log.Printf("Created directory: %s", path.Dir(symlink))
	return b.FS.Symlink(targetFile,symlink)
}

func (b Blog) getSimpleRepoPostFilePath(title string) string {
	return constructRepoPostFilePath(b.RepoPath,title)
}

type Article struct {
	Meta Metadata
	File io.Writer
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







