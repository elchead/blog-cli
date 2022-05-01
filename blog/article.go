package blog

import (
	"io"
	"path/filepath"
)

const articleDir = "/Blog"

type Article struct {
	Meta Metadata
	path string
	baseDir string
}

func NewArticleWithBaseDir(meta Metadata,baseDir string) *Article {
	return &Article{Meta:meta,baseDir:baseDir}
}

func (a Article) Title() string {
	return a.Meta.Title
}

func (a Article) Path() string { if a.path != "" { return a.path } else { return GetFilepath(a.Meta.Title,filepath.Join(a.baseDir,articleDir)) }  }

func (a Article) RepoFolder() string { return "posts" }

func (a Article) Write(file io.Writer) {
	io.WriteString(file,a.Meta.String())
}
