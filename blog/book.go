package blog

import (
	"io"
	"path/filepath"
)

const bookDir = "/Books"
const bookCategory = "Book-notes"

type Book struct {
	TemplateFile io.Reader
	Meta Metadata
	path string
	baseDir string
}

func NewBookWithBaseDir(meta Metadata,baseDir string) *Book {
	return &Book{Meta:meta,baseDir:baseDir}
}

func (b Book) Title() string {
	return b.Meta.Title
}

func (b Book) Path() string { if b.path != "" { return b.path } else {return GetFilepath(b.Meta.Title,filepath.Join(b.baseDir,bookDir)) } }

func (b Book) Write(bookFile io.Writer) {
	io.Copy(bookFile, b.TemplateFile)
}

func (b Book) RepoFolder() string { return "books" }
