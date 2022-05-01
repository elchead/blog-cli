package blog_test

import (
	"reflect"
	"testing"

	"github.com/elchead/blog-cli/blog"
	"github.com/stretchr/testify/assert"
)

func getType(post blog.Post) string {
	return reflect.TypeOf(post).String()
}

func TestCreatePost(t *testing.T) {
	postFactory := blog.PostFactory{BookTemplate: bookTemplate,BaseDir:baseDir}
	meta := blog.Metadata{Title: "Book title"}
	t.Run("letter", func(t *testing.T) {
		meta.Categories = []string{"Letters"}
		post, err := postFactory.NewPost(meta)
		assert.NoError(t, err)
		assert.Equal(t, "*blog.Letter", getType(post))
	})
	t.Run("book", func(t *testing.T) {
		meta.Categories = []string{"Book-notes"}
		post, err := postFactory.NewPost(meta)
		assert.NoError(t, err)
		assert.Equal(t, "*blog.Book", getType(post))
	})
	t.Run("article", func(t *testing.T) {
		meta.Categories = []string{"Programming"}
		post, err := postFactory.NewPost(meta)
		assert.NoError(t, err)
		assert.Equal(t, "*blog.Article", getType(post))
	})
}

// func TestFilePathFromPost(t *testing.T) {
// 	t.Run("letter", func(t *testing.T) {
// 		meta:= 
// 		post := Letter{"Letter title"}
// 		assert.Equal(t, "Blog/Book title.md", post.Path())
// 	})
// 	t.Run("book", func(t *testing.T) {
// 		post := NewBook("Book title", true)
// 		assert.Equal(t, "Books/Book title.md", post.Path())
// 	})
// 	t.Run("article", func(t *testing.T) {
// 		post := NewArticle("Book title", true)
// 		assert.Equal(t, "Blog/Book title.md", post.Path())
// 	})
// }
