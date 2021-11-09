package git

import (
	"fmt"
	"os/exec"

	"github.com/pkg/errors"
)

type Repo struct {
	RepoPath string
}

func (r *Repo) Pull() error {
	return r.execCommand("pull")
}
func (r *Repo) StageAll() error {
	return r.execCommand("add",".")
}

func (r *Repo) Commit(title string) error {
	return r.execCommand("commit","-m",fmt.Sprintf("New post: %s",title))
}

func (r *Repo) Push() error {
	return r.execCommand("push")
}

func (r *Repo) execCommand(args ...string) error {
	cmd := exec.Command("git", args...)
	cmd.Dir = r.RepoPath
	output, _ := cmd.CombinedOutput()
	if std:=string(output); std!= "" { 
		fmt.Println(std) 
	}
	if err := cmd.Run(); err != nil {
		return errors.Wrapf(err, "Git %s failed: %v", args[0], err)
	}
	return nil
}
