package blog

import (
	"fmt"
	"io"
	"path"
)

func GetFilepath(articleTitle,folderPath string) string {
	return path.Join(folderPath,articleTitle+".md")
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
}

func (b Blog) WritePost(metadata Metadata,file io.Writer) {
	io.WriteString(file,metadata.String())
}
