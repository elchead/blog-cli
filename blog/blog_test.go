package blog_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path"
	"strings"
	"testing"

	"github.com/elchead/blog-cli/blog"
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


func TestCreateFile(t *testing.T) {
	t.Run("save file",func(t *testing.T) {
		appFS := afero.NewMemMapFs()
		title := "title"
		path := "blog-cli"
		appFS.Create(blog.GetFilepath(title,path))
		exists, _ := afero.Exists(appFS,path+"/title.md")
		assert.Equal(t,true,exists)
	})
	t.Run("non existent file is not saved", func(t *testing.T) {
		title := "title"
		path := "blog-cli"
		appFS := afero.NewMemMapFs()
		appFS.Create(blog.GetFilepath(title,path))
		exists, _ := afero.Exists(appFS,path+"nofile.md")
		assert.Equal(t,false,exists)
	})
}


func TestBlog(t *testing.T) {
	mockedFs := afero.NewMemMapFs()
	fakeFs := &FakeSymLinker{Fs: mockedFs}
	sut := blog.Blog{RepoPath: "/repo",WritingDir:"/writing",FS:fakeFs}
	meta := blog.Metadata{Title: "Learning is great - Doing is better", Categories : []string{"Thoughts"}, Date: "2021-11-04"}
	t.Run("file is created in expected directory", func(t *testing.T){
		_,err := sut.DraftPost(meta) 
		assert.NoError(t, err)
		_, err = mockedFs.Open(path.Join(sut.WritingDir,meta.Title+".md"))
		assert.NoError(t,err)
	})
	t.Run("link article in repo",func(t *testing.T){
		article,_ := sut.DraftPost(meta) 
		err := sut.LinkInRepo(article)
		assert.NoError(t,err)
		wantedDirName := "learning-is-great"
		wantedSymlink := path.Join(sut.RepoPath,"content","posts",wantedDirName,"index.en.md")
		_, err = mockedFs.Open(wantedSymlink)
		assert.NoError(t,err)

	})
	t.Run("draft book",func(t *testing.T){
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
}

func TestArticle(t *testing.T){
	meta := blog.Metadata{Title: "Learning is great - Doing is better", Categories : []string{"Thoughts"}, Date: "2021-11-04"}
	sut := blog.Article{Meta: meta}
	t.Run("write meta to io.Writer", func(t *testing.T) {
		var file bytes.Buffer
		sut.Write(&file)
		assert.Equal(t,meta.String(),file.String())
	})
}

func TestFakeSymLink(t *testing.T){
	mockedFs := afero.NewMemMapFs()
	sut := &FakeSymLinker{Fs:mockedFs}
	// files don't exist!
	assert.Error(t,sut.Symlink("/file/write.md","link.md"))
}

type FakeSymLinker struct {
	afero.Fs
	TargetFile string
	CreatedSymlink string
}

func assertDirExists(fs afero.Fs, path string) error {
	exists, err := afero.Exists(fs, path)
	if !exists || err != nil {
		return fmt.Errorf("directory (%s) does not exist: %w",path,err)
	}
	return nil
}

func (f* FakeSymLinker) Symlink(target, link string) error {
	if err := assertDirExists(f,path.Dir(link)); err != nil {
		return err
	}
	
	if err:= assertDirExists(f,path.Dir(target)); err != nil {
		return err
	}
	f.TargetFile = target
	f.CreatedSymlink = link
	f.Create(link)
	return nil
}
