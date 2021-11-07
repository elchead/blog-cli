package blog

import "io"

type Book struct {
	TemplateFile io.Reader
}

func (b Book) Write(bookFile io.Writer) error {
	_,err := io.Copy(bookFile, b.TemplateFile)
	return err
}
