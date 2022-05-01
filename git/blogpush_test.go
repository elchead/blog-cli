package git

import (
	"testing"

	"github.com/elchead/blog-cli/blog"
	"github.com/stretchr/testify/assert"
)

type fakeRepo struct {
	repoPath string
	callOrder []string
}

func (r *fakeRepo) StageAll() error {
	r.callOrder = append(r.callOrder, "StageAll")
	return nil
}

func (r *fakeRepo) Commit(title string) error {
	r.callOrder = append(r.callOrder, "Commit")
	return nil
}

func (r *fakeRepo) Pull() error {
	r.callOrder = append(r.callOrder, "Pull")
	return nil
}

func (r *fakeRepo) Push() error {
	r.callOrder = append(r.callOrder, "Push")
	return nil
}

func (r *fakeRepo) RepoPath() string {
	return r.repoPath
}


func TestPushOrder(t *testing.T) {
	fakeRepo := &fakeRepo{repoPath:""}
	sut := BlogPush{fakeRepo}
	sut.Push(blog.NewArticleWithBaseDir(blog.Metadata{},"/"))
	assert.Equal(t, []string{"StageAll", "Commit", "Pull", "Push"}, fakeRepo.callOrder)
}
