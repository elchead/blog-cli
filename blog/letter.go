package blog

import "io"

type Letter struct {
	TemplateFile io.Reader
	Meta Metadata
	path string
}

// constructor ensures that path is always provided for safety
func NewLetter(meta Metadata,path string) *Letter {
	return &Letter{Meta:meta,path:path}
}

func (b Letter) Title() string {
	return b.Meta.Title
}

func (b Letter) Path() string { return b.path }

func (a Letter) Write(file io.Writer) {
	io.WriteString(file,a.Meta.String())
}

func (b Letter) RepoFolder() string { return "letters" }

