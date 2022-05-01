package blog

import (
	"io"
	"path/filepath"
)

const letterCategory = "Letters"
const letterDir = "/Letters"


type Letter struct {
	TemplateFile io.Reader
	Meta Metadata
	baseDir string
}

func NewLetterWithBaseDir(meta Metadata, baseDir string) *Letter {
	return &Letter{Meta:meta,baseDir:baseDir}
}

func (b Letter) Title() string {
	return b.Meta.Title
}

func (b Letter) Path() string { return GetFilepath(b.Meta.Title,filepath.Join(b.baseDir,letterDir))}

func (a Letter) Write(file io.Writer) {
	io.WriteString(file,a.Meta.String())
}

func (b Letter) RepoFolder() string { return "letters" }

