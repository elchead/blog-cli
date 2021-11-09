package blog

import "io"

type Article struct {
	Meta Metadata
	File io.Writer
	path string
}

// constructor ensures that path is always provided for safety
func NewArticle(meta Metadata,path string) *Article {
	return &Article{Meta:meta,path:path}
}

func (a Article) Title() string {
	return a.Meta.Title
}

func (a Article) Path() string { return a.path }

func (a Article) RepoFolder() string { return "posts" }

func (a Article) Write(file io.Writer) {
	io.WriteString(file,a.Meta.String())
}
