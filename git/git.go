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
	cmd := exec.Command("git","pull")
	cmd.Dir = r.RepoPath
	output, _ := cmd.CombinedOutput()
	fmt.Println(string(output))
	if err := cmd.Run(); err != nil {
		return errors.Wrapf(err, "Git pull failed: %v", err)
	}
	return nil
}
func (r *Repo) StageAll() error {
	cmd := exec.Command("git","add",".")
	cmd.Dir = r.RepoPath
	output, _ := cmd.CombinedOutput()
	fmt.Println(string(output))
	if err := cmd.Run(); err != nil {
		return errors.Wrapf(err, "Git stage all failed: %v", err)
	}
	return nil
}

func (r *Repo) Commit(title string) error {
	cmd := exec.Command("git","commit","-m",fmt.Sprintf("New post: %s",title))
	cmd.Dir = r.RepoPath
	output, _ := cmd.CombinedOutput()
	fmt.Println(string(output))
	if err := cmd.Run(); err != nil {
		return errors.Wrapf(err, "Git commit failed: %v", err)
	}
	return nil
}

func (r *Repo) Push() error {
	cmd := exec.Command("git","push") //("git","push","--set-upstream","origin","test") when first time!!
	cmd.Dir = r.RepoPath
	output, _ := cmd.CombinedOutput()
	fmt.Println(string(output))
	if err := cmd.Run(); err != nil {
		return errors.Wrapf(err, "Git push failed: %v", err)
	}
	return nil
}

// type BlogGit struct {}

// func (g *BlogGit) Publish(post) error {}
