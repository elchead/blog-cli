package blog

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"
)

func GetFilepath(articleTitle,folderPath string) string {
	return path.Join(folderPath,articleTitle+".md")
}

type FsSymLinker interface {
	Symlink(target, symlink string) error
	MkdirAll(path string, perm os.FileMode) error 
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
}

func (b Blog) WritePost(metadata Metadata,file io.Writer) {
	io.WriteString(file,metadata.String())
}

func (b Blog) getRepoPostFilePath(meta Metadata) string {
	return path.Join(b.RepoPath,"content","posts",meta.Title,"index.en.md") // TODO shorten directory name of article
}

func (b Blog) CreatePostInRepo(fsys FsSymLinker,meta Metadata,targetFile string) error {
	symlink := b.getRepoPostFilePath(meta)
	err := fsys.MkdirAll(path.Dir(symlink),0777)
	if err != nil {
		return fmt.Errorf("could not create directory: %w", err)
	}
	err = fsys.MkdirAll(path.Dir(targetFile),0777)
	if err != nil {
		return fmt.Errorf("could not create directory: %w", err)
	}
	log.Printf("Created directory: %s", path.Dir(symlink))
	return fsys.Symlink(targetFile,symlink)
}




