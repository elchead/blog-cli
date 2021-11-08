package blog

import "io"

type Article struct {
	Meta Metadata
	File io.Writer
	Path string
}

func (a Article) Title() string {
	return a.Meta.Title
}

func (a Article) RepoFolder() string { return "posts" }

func (b Article) Write(file io.Writer) {
	io.WriteString(file,b.Meta.String())
}
