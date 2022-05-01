package blog

import "io"

const bookDir = obsidianVault +"/Books"

type Book struct {
	TemplateFile io.Reader
	Meta Metadata
	path string
}

// constructor ensures that path is always provided for safety
func NewBookWithPath(meta Metadata,path string) *Book {
	return &Book{Meta:meta,path:path}
}

func NewBook(meta Metadata) *Book {
	return &Book{Meta:meta}
}

func (b Book) Title() string {
	return b.Meta.Title
}

func (b Book) Path() string { return GetFilepath(b.Meta.Title,bookDir) }

func (b Book) Write(bookFile io.Writer) {
	io.Copy(bookFile, b.TemplateFile)
}

func (b Book) RepoFolder() string { return "books" }
