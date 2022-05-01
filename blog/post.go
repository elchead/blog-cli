package blog

import (
	"errors"
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


type PostFactory struct {
	BookTemplate io.Reader
}

func (f PostFactory) NewPost(meta Metadata,baseDir string) (post Post,err error) {
	if f.BookTemplate == nil {
		return nil, errors.New("Book template not defined")
	}
	if contains(meta.Categories, letterCategory) && contains(meta.Categories, bookCategory) {
		return nil, fmt.Errorf("post category ambiguous. Found both letter and book")
	}
	switch {
	case contains(meta.Categories,bookCategory):
		book := NewBookWithBaseDir(meta,baseDir)
		book.TemplateFile = f.BookTemplate
		return book, nil
	case contains(meta.Categories,letterCategory):
		return NewLetter(meta),nil
	default:
		return NewArticleWithBaseDir(meta,baseDir),nil
	}	
}

func NewPost(meta Metadata) (post Post,err error) {
	return NewPostWithBaseDir(meta,obsidianVault)
}

func NewPostWithBaseDir(meta Metadata,baseDir string) (post Post,err error) {
	if contains(meta.Categories, letterCategory) && contains(meta.Categories, bookCategory) {
		return nil, fmt.Errorf("post category ambiguous. Found both letter and book")
	}
	switch {
	case contains(meta.Categories,bookCategory):
		return NewBookWithBaseDir(meta,baseDir),nil
	case contains(meta.Categories,letterCategory):
		return NewLetter(meta),nil
	default:
		return NewArticleWithBaseDir(meta,baseDir),nil
	}
}
