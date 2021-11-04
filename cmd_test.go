package blog_test

import (
	"testing"

	"github.com/elchead/blog-cli"
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
