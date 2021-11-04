package blog_test

import (
	"bytes"
	"testing"

	"github.com/elchead/blog-cli"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)
func TestAddMetadata(t *testing.T) {
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
		sut := blog.File{Title: "title",Path:"blog-cli"}
		appFS := afero.NewMemMapFs()
		sut.Create(appFS)
		exists, _ := afero.Exists(appFS,sut.Path+"/title.md")
		assert.Equal(t,true,exists)
	})
	t.Run("non existent file is not saved", func(t *testing.T) {
		sut := blog.File{Title: "title",Path:"blog-cli"}
		appFS := afero.NewMemMapFs()
		sut.Create(appFS)
		exists, _ := afero.Exists(appFS,sut.Path+"nofile.md")
		assert.Equal(t,false,exists)
	})
}

func TestBlog(t *testing.T){
	sut := blog.Blog{}
	meta := blog.Metadata{Title: "title", Categories : []string{"Thoughts"}, Date: "2021-11-04"}
	var file bytes.Buffer
	sut.WritePost(meta,&file)
	assert.Equal(t,meta.String(),file.String())
}
