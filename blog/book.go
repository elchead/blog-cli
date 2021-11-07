package blog

import "io"

type Book struct {
	TemplateFile io.Reader
	Meta Metadata
}

func (b Book) Title() string {
	return b.Meta.Title
}

func (b Book) Write(bookFile io.Writer) {
	io.Copy(bookFile, b.TemplateFile)
	//return err
}
