package blog

import "io"

type Book struct {
	TemplateFile io.Reader
	Meta Metadata
	path string
}

// constructor ensures that path is always provided for safety
func NewBook(meta Metadata,path string) *Book {
	return &Book{Meta:meta,path:path}
}

func (b Book) Title() string {
	return b.Meta.Title
}

func (b Book) Path() string { return b.path }

func (b Book) Write(bookFile io.Writer) {
	io.Copy(bookFile, b.TemplateFile)
}

func (b Book) RepoFolder() string { return "books" }
