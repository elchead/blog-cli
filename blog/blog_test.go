package blog_test

import (
	"io/ioutil"
	"path"
	"strings"
	"testing"

	"github.com/elchead/blog-cli/blog"
	"github.com/elchead/blog-cli/fs"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)
func TestMetadata(t *testing.T) {
	sut := blog.Metadata{Title: "title", Categories : []string{"Thoughts"}, Date: "2021-11-04"}
	want := `---
title: title
categories: [Thoughts]
date: 2021-11-04
---`
	assert.Equal(t,want,sut.String())
}

func TestBlog(t *testing.T) {
	mockedFs := afero.NewMemMapFs()
	fakeFs := &fs.FakeSymLinker{Fs: mockedFs}
	sut := blog.Blog{RepoPath: "/repo",WritingDir:"/writing",FS:fakeFs}
	meta := blog.Metadata{Title: "Learning is great - Doing is better", Categories : []string{"Thoughts"}, Date: "2021-11-04"}
	t.Run("article is created in expected directory", func(t *testing.T){
		_,err := sut.DraftArticle(meta) 
		assert.NoError(t, err)
		_, err = mockedFs.Open(path.Join(sut.WritingDir,meta.Title+".md"))
		assert.NoError(t,err)
	})
	t.Run("link article in repo",func(t *testing.T){
		article,_ := sut.DraftArticle(meta) 
		err := sut.LinkInRepo(article)
		assert.NoError(t,err)
		wantedDirName := "learning-is-great"
		wantedSymlink := path.Join(sut.RepoPath,"content","posts",wantedDirName,"index.en.md")
		_, err = mockedFs.Open(wantedSymlink)
		assert.NoError(t,err)

	})
	t.Run("book is created in expected directory",func(t *testing.T){
		sut.BookTemplate = strings.NewReader("---template---")
		sut.BookDir = "/writing/Books"
		_,err := sut.DraftBook(blog.Metadata{Title: "Alice"}) 
		assert.NoError(t,err)
		file, err := mockedFs.Open("/writing/Books/Alice.md")
		assert.NoError(t,err)
		content,err := ioutil.ReadAll(file)
		assert.NoError(t,err)
		assert.Equal(t,"---template---",string(content))
	})
	t.Run("link book in repo",func(t *testing.T){
		sut.BookTemplate = strings.NewReader("---template---")
		sut.BookDir = "/writing/Books"
		book,_ := sut.DraftBook(blog.Metadata{Title: "Alice"}) 
		sut.LinkInRepo(book)
		wantedSymlink := path.Join(sut.RepoPath,"content","books","alice","index.en.md")
		_, err := mockedFs.Open(wantedSymlink)
		assert.NoError(t,err)
		assert.Equal(t,sut.BookDir+"/Alice.md",fakeFs.TargetFile)
	})
}
