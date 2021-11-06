package blog

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"strings"
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

func constructDirNameFromTitle(title string) string {
	lowerCase := strings.ToLower(title)
	cutAfterDash := strings.Split(lowerCase," - ")[0]
	noSpaces := strings.Replace(cutAfterDash, " ","-",-1)
	return noSpaces
}

func constructRepoPostFilePath(repoPath ,dirName string) string {
	return path.Join(repoPath,"content","posts",constructDirNameFromTitle(dirName),"index.en.md")
}

func (b Blog) getSimpleRepoPostFilePath(meta Metadata) string {
	return constructRepoPostFilePath(b.RepoPath,meta.Title)
}

func (b Blog) CreatePostInRepo(fsys FsSymLinker,meta Metadata,targetFile string) error {
	symlink := b.getSimpleRepoPostFilePath(meta)
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


func (b Blog) DraftPost(fsys FsSymLinker, meta Metadata) {
	return
}




