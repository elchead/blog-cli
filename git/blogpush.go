package git

import (
	"log"

	"github.com/elchead/blog-cli/blog"
	"github.com/elchead/blog-cli/fs"
)

type Repoer interface {
	StageAll() error
	Commit(title string) error
	Pull() error
	Push() error
	RepoPath() string
}

type BlogPush struct { 
	repo Repoer
}

func NewBlogPush(repoPath string) *BlogPush {
	return &BlogPush{repo: &Repo{repoPath:repoPath}}
}

func (p *BlogPush) Push(post blog.Post) error {
	log.Print("Preparing Git repo for publishing...\n")
	symlink := blog.ConstructRepoPostFilePath(p.repo.RepoPath(),post.RepoFolder(),post.Title())
	fs.MakeHardlink(symlink)
	p.repo.StageAll()
	p.repo.Commit(post.Title())
	p.repo.Pull() // TODO ignore error: exec: already started
	p.repo.Push()
	return nil
}

func (p *BlogPush) PushChanges() error {
	log.Print("Pushing changes...\n")
	p.repo.StageAll()
	p.repo.Commit("Push changes from CLI cmd")
	p.repo.Pull() // TODO ignore error: exec: already started
	p.repo.Push()
	return nil
}


