package blog

import "io"

type Book struct {
	TemplateFile io.Reader
	Meta Metadata
	Path_ string
}

func (b Book) Title() string {
	return b.Meta.Title
}

func (b Book) Path() string { return b.Path_ }

func (b Book) Write(bookFile io.Writer) {
	io.Copy(bookFile, b.TemplateFile)
}

func (b Book) RepoFolder() string { return "books" }
