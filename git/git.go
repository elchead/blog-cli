package git

import (
	"fmt"
	"os/exec"

	"github.com/pkg/errors"
)

type Repo struct {
	repoPath string
}

func (r *Repo) Pull() error {
	return r.execCommand("pull")
}
func (r *Repo) StageAll() error {
	return r.execCommand("add",`.`)
}

func (r *Repo) Commit(title string) error {
	return r.execCommand("commit","-m",fmt.Sprintf("New post: %s",title))
}

func (r *Repo) Push() error {
	return r.execCommand("push")
}

func (r *Repo) RepoPath() string {
	return r.repoPath
}

func (r *Repo) execCommand(args ...string) error {
	cmd := exec.Command("git", args...)
	cmd.Dir = r.repoPath
	output,err := cmd.CombinedOutput()
	if std:=string(output); std!= "" { 
		fmt.Println(std) 
	}
	if err != nil {
		return errors.Wrapf(err, "Git %s failed: %v", args[0], err)
	}
	return nil
}
