package blog

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func getType(post Post) string {
	return reflect.TypeOf(post).String()
}

func TestCreatePost(t *testing.T) {
	meta := Metadata{Title: "Book title"}
	t.Run("letter", func(t *testing.T) {
		meta.Categories = []string{"Letter"}
		post, err := NewPost(meta)
		assert.NoError(t, err)
		assert.Equal(t, "*blog.Letter", getType(post))
	})
	t.Run("book", func(t *testing.T) {
		meta.Categories = []string{"Book-notes"}
		post, err := NewPost(meta)
		assert.NoError(t, err)
		assert.Equal(t, "*blog.Book", getType(post))
	})
	t.Run("article", func(t *testing.T) {
		meta.Categories = []string{"Programming"}
		post, err := NewPost(meta)
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
