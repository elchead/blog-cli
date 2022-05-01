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
		post, err := createPost(meta, "/")
		assert.NoError(t, err)
		assert.Equal(t, "*blog.Letter", getType(post))
	})
	t.Run("book", func(t *testing.T) {
		meta.Categories = []string{"Book-notes"}
		post, err := createPost(meta, "/")
		assert.NoError(t, err)
		assert.Equal(t, "*blog.Book", getType(post))
	})
	t.Run("article", func(t *testing.T) {
		meta.Categories = []string{"Programming"}
		post, err := createPost(meta, "/")
		assert.NoError(t, err)
		assert.Equal(t, "*blog.Article", getType(post))
	})

}
