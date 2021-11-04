package blog

import (
	"fmt"
	"io"
	"path"

	"github.com/spf13/afero"
)

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

type File struct {
	Title string
	Path string
}	


func (f File) Filepath() string {
	return path.Join(f.Path,f.Title+".md")
}

func (f File) Create(fsys afero.Fs) (afero.File, error) {
	return fsys.Create(f.Filepath())
}

type Blog struct {
}

func (b Blog) WritePost(metadata Metadata,file io.Writer) {
	io.WriteString(file,metadata.String())
}
