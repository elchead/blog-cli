package blog

import "io"

type Article struct {
	Meta Metadata
	File io.Writer
	Path_ string
}

func (a Article) Title() string {
	return a.Meta.Title
}

func (a Article) Path() string { return a.Path_ }

func (a Article) RepoFolder() string { return "posts" }

func (a Article) Write(file io.Writer) {
	io.WriteString(file,a.Meta.String())
}
