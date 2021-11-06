package blog_test

import (
	"bytes"
	"fmt"
	"os"
	"path"
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

func TestBlog(t *testing.T){
	sut := blog.Blog{RepoPath: "/repo"}
	meta := blog.Metadata{Title: "Learning is great - Doing is better", Categories : []string{"Thoughts"}, Date: "2021-11-04"}
	writingPath := blog.GetFilepath(meta.Title,"/writing")
	t.Run("write meta to io.Writer", func(t *testing.T) {
		var file bytes.Buffer
		sut.WritePost(meta,&file)
		assert.Equal(t,meta.String(),file.String())
	})
	t.Run("create repo skeleton with shortened directory name", func(t *testing.T){
		mockedFs := afero.NewMemMapFs()
		fakeFs := &FakeSymLinker{fs: mockedFs,t: t}

		err := sut.CreatePostInRepo(fakeFs,meta.Title,writingPath)
		assert.NoError(t,err)
		wantedDirName := "learning-is-great"
		wantedSymlink := path.Join(sut.RepoPath,"content","posts",wantedDirName,"index.en.md")
		_, err = mockedFs.Open(wantedSymlink)
		assert.NoError(t,err)
	})
	// t.Run("draft post in writing path without git repo",func(t *testing.T) {
	// 	mockedFs := afero.NewMemMapFs()
	// 	fakeFs := &FakeSymLinker{fs: mockedFs,t: t}
	// 	sut.DraftPost(fakeFs,meta)
	// 	mockedFs.Open()

	// })
}

func TestFakeSymLink(t *testing.T){
	mockedFs := afero.NewMemMapFs()
	sut := &FakeSymLinker{fs: mockedFs,t: t}	
	// files don't exist!
	assert.Error(t,sut.Symlink("/file/write.md","link.md"))
}

type FakeSymLinker struct {
	TargetFile string
	CreatedSymlink string
	fs afero.Fs
	t testing.TB
}

func assertDirExists(fs afero.Fs, path string) error {
	exists, err := afero.Exists(fs, path)
	if !exists || err != nil {
		return fmt.Errorf("directory (%s) does not exist: %w",path,err)
	}
	return nil
}

func (f* FakeSymLinker) Symlink(target, link string) error {
	if err := assertDirExists(f.fs,path.Dir(link)); err != nil {
		return err
	}
	
	if err:= assertDirExists(f.fs,path.Dir(target)); err != nil {
		return err
	}
	f.TargetFile = target
	f.CreatedSymlink = link
	f.fs.Create(link)
	return nil
}

func (f* FakeSymLinker) MkdirAll(path string, perm os.FileMode) error {
	return f.fs.MkdirAll(path, perm)
}

