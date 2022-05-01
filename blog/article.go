package blog

import "io"

const articleDir = obsidianVault +"/Blog"

type Article struct {
	Meta Metadata
	File io.Writer
	path string
}

// constructor ensures that path is always provided for safety
func NewArticleWithPath(meta Metadata,path string) *Article {
	return &Article{Meta:meta,path:path}
}

func NewArticle(meta Metadata) *Article {
	return &Article{Meta:meta}
}

func (a Article) Title() string {
	return a.Meta.Title
}

func (a Article) Path() string { if a.path != "" { return a.path } else { return GetFilepath(a.Meta.Title,articleDir) }  }

func (a Article) RepoFolder() string { return "posts" }

func (a Article) Write(file io.Writer) {
	io.WriteString(file,a.Meta.String())
}
