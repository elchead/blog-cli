package main

import (
	"os"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type mockFn struct{
	calledArg string
}
func (m *mockFn) call(path string){
	m.calledArg = path
}


func TestPathAndFilenameExtraction(t *testing.T) {
	path := "/Users/a/Blog/post_title.md"
	assert.Equal(t, "Blog/post_title.md",GetVaultPath(path))

}
func TestPushToReadwise(t *testing.T) {
	reader := strings.NewReader("y!\n")
	post := newPost("post_title",true)
	mockFn := mockFn{}
	AskToPublishToReadwise(reader,post,mockFn.call)
	assert.Equal(t,repoDir+"/content/books/post_title",mockFn.calledArg)
}

func TestAskToPublish(t *testing.T) {
	t.Run("input yes",func(t *testing.T) {
		reader := strings.NewReader("y!\n")
		assert.True(t,okToPublish(reader))
	})
	t.Run("input no",func(t *testing.T) {
		reader := strings.NewReader("n!\n")
		assert.False(t,okToPublish(reader))
	})
	t.Run("input wrong yes",func(t *testing.T) {
		reader := strings.NewReader("y\n")
		assert.False(t,okToPublish(reader))
	})
}

// not needed at the moment
func TestExitOfRenderRoutine(t *testing.T) {
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	startGoRoutine(sigs,done)
	sigs <- syscall.SIGINT
	assert.True(t,checkIfExited(500*time.Millisecond,done))
}

func checkIfExited(timeout time.Duration, done chan bool) bool {
	for {
		select {
		    case <-done:
			return true
		    case <-time.After(timeout):
			return false
		}
	    }
}
