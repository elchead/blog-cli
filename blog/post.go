package blog

import (
	"fmt"
	"io"
)


const letterCategory = "Letter"
const bookCategory = "Book-notes"

type Post interface {
	Title() string
	Write(file io.Writer)
	RepoFolder() string
	Path() string
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func createPost(meta Metadata, path string) (post Post,err error) {
	if contains(meta.Categories, letterCategory) && contains(meta.Categories, bookCategory) {
		return nil, fmt.Errorf("post category ambiguous. Found both letter and book")
	}
	switch {
	case contains(meta.Categories,bookCategory):
		return NewBook(meta,path),nil
	case contains(meta.Categories,letterCategory):
		return NewLetter(meta,path),nil
	default:
		return NewArticle(meta,path),nil
	}
}
