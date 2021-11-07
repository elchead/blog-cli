package blog

import "io"

type Book struct {
 	Article
}

func (b Book) CreateNote(template io.Reader,bookFile io.Writer) error {
	_,err := io.Copy(bookFile,template)
	return err
}
