package main

import (
	"log"
	"strings"
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

func TestNewMetadata(t *testing.T){
	t.Run("letter",func(t *testing.T) {
		meta := newMetadata("test",false,true)
		assert.Equal(t,"Letters",meta.Categories[0])
		assert.Equal(t,"test",meta.Title)
	})
	t.Run("book",func(t *testing.T) {
		meta := newMetadata("test",true,false)
		assert.Equal(t,"Book-notes",meta.Categories[0])
		assert.Equal(t,"test",meta.Title)
	})
	t.Run("article",func(t *testing.T) {
		input := strings.NewReader("Thoughts\n")
		meta := newMetadataFrom("Test title",false,false,input)
		assert.Equal(t,"Thoughts",meta.Categories[0])
		assert.Equal(t,"Test title",meta.Title)
	})
}

func TestPathAndFilenameExtraction(t *testing.T) {
	path := "/Users/a/Blog/post_title.md"
	assert.Equal(t, "Blog/post_title.md",GetVaultPath(path))

}
func TestPushToReadwise(t *testing.T) {
	reader := strings.NewReader("y!\n")
	meta := newMetadata("post title",true,false)
	post,err := postFactory.NewPost(meta)
	if err != nil {
		log.Fatal(err)
	}
	mockFn := mockFn{}
	AskToPublishToReadwise(reader,post,mockFn.call)
	assert.Equal(t,repoDir+"/content/books/post-title",mockFn.calledArg)
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

// // not needed at the moment
// func TestExitOfRenderRoutine(t *testing.T) {
// 	sigs := make(chan os.Signal, 1)
// 	done := make(chan bool, 1)
// 	startGoRoutine(sigs,done)
// 	sigs <- syscall.SIGINT
// 	assert.True(t,checkIfExited(500*time.Millisecond,done))
// }

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
