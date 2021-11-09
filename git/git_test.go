package git_test

import (
	"testing"

	"github.com/elchead/blog-cli/git"
	"github.com/stretchr/testify/assert"
)

func TestGit(t *testing.T) {
	// directory := "/Users/adria/Programming/elchead.github.io"
	// username := "elchead"
	// password := "super-secret-Personal-access-token"
	// repo := "github.build.company.com/org/repo-name.git"
	// url := fmt.Sprintf("https://%s:%s@%s", username, password, repo)
	// r, err := git.PlainOpen(directory)
	// assert.NoError(t, err)
	// w, err := r.Worktree()
// 	pull := exec.Command("git","add",".")
// 	pull.Dir = "/Users/adria/Programming/elchead.github.io"
// 	err := pull.Run()
// 	commit := exec.Command("git","commit","-m","new post exec")
// 	commit.Dir = "/Users/adria/Programming/elchead.github.io"
// err = commit.Run()
// 	assert.NoError(t, err)
	// exec.Command()
	// w.Pull(&git.PullOptions{})
	// Print the latest commit that was just pulled
	// ref, err := r.Head()
	//  r.CommitObject(ref.Hash())
	//  // fmt.Println(commit)
	//  assert.Error(t, err)
	//  // fmt.Printf(":")
	//  err = w.AddGlob("content/posts/*/index.en.md")
	//  assert.NoError(t, err)
	//  _, err = w.Commit("go git post", &git.CommitOptions{})
			 
	// err = r.Push(&git.PushOptions{})
	// assert.NoError(t,err)
	
}


func TestPostPublish(t *testing.T) {
	repo := git.Repo{RepoPath: "/Users/adria/Programming/elchead.github.io"}
	err := repo.Pull()
	assert.NoError(t, err)
	err = repo.StageAll()
	assert.NoError(t, err)
	err = repo.Commit("hioo")
	assert.NoError(t, err)
	err = repo.Push()
	assert.NoError(t, err)
	


}


// func TestBlogGit(t *testing.T){
// 	repo := git.BlogGit{}
// 	repo.Publish(post)
// }
