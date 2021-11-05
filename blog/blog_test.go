package blog_test

import (
	"bytes"
	"log"
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
	sut := blog.Blog{RepoPath: "/blog"}
	meta := blog.Metadata{Title: "title", Categories : []string{"Thoughts"}, Date: "2021-11-04"}
	fpath := blog.GetFilepath(meta.Title,sut.RepoPath)
	t.Run("write meta to file", func(t *testing.T) {
		var file bytes.Buffer
		sut.WritePost(meta,&file)
		assert.Equal(t,meta.String(),file.String())
	})
	t.Run("create repo skeleton", func(t *testing.T){
		mockedFs := afero.NewMemMapFs()
		fakeFs := &FakeSymLinker{fs: mockedFs}
		err := sut.CreatePost(fakeFs,meta,fpath)
		assert.NoError(t,err)
		wantedSymlink := path.Join(sut.RepoPath,"content","posts",meta.Title,"index.en.md")
		assert.Equal(t,wantedSymlink,fakeFs.CreatedSymlink)
		assert.Equal(t,fpath,fakeFs.TargetFile)
	})
}

type FakeSymLinker struct {
	TargetFile string
	CreatedSymlink string
	fs afero.Fs
}

func (f* FakeSymLinker) Symlink(target, link string) error {
	exists,err := afero.Exists(f.fs,path.Dir(target))
	if err!=nil || !exists {
		log.Fatalf("Directory not found: %v", err)
	}
	f.TargetFile = target
	f.CreatedSymlink = link
	return nil
}

func (f* FakeSymLinker) MkdirAll(path string, perm os.FileMode) error {
	return f.fs.MkdirAll(path, perm)
}

