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
resources:
  - name: "featured-image"
    src: "cover.jpg"
date: 2021-11-04
---`
	assert.Equal(t,want,sut.String())
}

var bookTemplate = strings.NewReader("---template---")
var postFactory = blog.PostFactory{BookTemplate: bookTemplate}
func TestBlog(t *testing.T) {
	
	blog.PpostFactory.BookTemplate = bookTemplate
	mockedFs := afero.NewMemMapFs()
	fakeFs := &fs.FakeSymLinker{Fs: mockedFs}
	sut := blog.BlogWriter{RepoPath: "/repo",WritingDir:"/writing",FS:fakeFs,BookTemplate:bookTemplate,BookDir:"/writing/Books",LetterDir:"/writing/Letters"}
	meta := blog.Metadata{Title: "Learning is great - Doing is better", Categories : []string{"Thoughts"}, Date: "2021-11-04"}
	t.Run("article is created in expected directory", func(t *testing.T){
		post, err := postFactory.NewPost(meta,sut.WritingDir)
		assert.NoError(t, err)
		err = sut.WritePost(post)
		assert.NoError(t, err)
		_, err = mockedFs.Open(path.Join(sut.WritingDir,"/Blog",meta.Title+".md"))
		assert.NoError(t,err)
	})
	t.Run("link article in repo",func(t *testing.T){
		article, err := postFactory.NewPost(meta,sut.WritingDir)
		assert.NoError(t,err)
		err = sut.LinkInRepo(article)
		assert.NoError(t,err)
		wantedDirName := "learning-is-great"
		wantedSymlink := path.Join(sut.RepoPath,"content","posts",wantedDirName,"index.en.md")
		_, err = mockedFs.Open(wantedSymlink)
		assert.NoError(t,err)

	})
	t.Run("book is created in expected directory",func(t *testing.T){
		meta := blog.Metadata{Title: "Alice",Categories : []string{"Book-notes"}, Date: "2021-11-04"} 
		post, err := postFactory.NewPost(meta,sut.WritingDir)	
		assert.NoError(t,err)
		err = sut.WritePost(post)
		assert.NoError(t, err)

		file, err := mockedFs.Open("/writing/Books/Alice.md")
		assert.NoError(t,err)
		content,err := ioutil.ReadAll(file)
		assert.NoError(t,err)
		assert.Equal(t,"---template---",string(content))
	})
	t.Run("link book in repo",func(t *testing.T){
		meta := blog.Metadata{Title: "Alice",Categories : []string{"Book-notes"}, Date: "2021-11-04"} 
		post, err := postFactory.NewPost(meta,sut.WritingDir)		
		assert.NoError(t, err)
		err = sut.LinkInRepo(post)
		assert.NoError(t, err)

		wantedSymlink := path.Join(sut.RepoPath,"content","books","alice","index.en.md")
		_, err = mockedFs.Open(wantedSymlink)
		assert.NoError(t,err)
		assert.Equal(t,sut.BookDir+"/Alice.md",fakeFs.TargetFile)
	})
	t.Run("add image to post",func(t *testing.T) {
		article, err := postFactory.NewPost(meta,sut.WritingDir)		
		assert.NoError(t, err)
		err = sut.LinkInRepo(article)
		assert.NoError(t,err)
		img := strings.NewReader("img")
		sut.AddMedia(article,img,"img.txt")
		
		wantedLink := path.Join(sut.RepoPath,"content","posts","learning-is-great","img.txt")
		_, err = mockedFs.Open(wantedLink)
		assert.NoError(t,err)
	})
}

func TestLink(t *testing.T) {
	post := blog.NewArticleWithPath(blog.Metadata{Title:"Examples are good"},"")
	assert.Equal(t,"posts/examples-are-good",blog.ConstructPostLink(post))
}
