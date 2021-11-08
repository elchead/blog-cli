package blog_test

import (
	"bytes"
	"testing"

	"github.com/elchead/blog-cli/blog"
	"github.com/stretchr/testify/assert"
)

func TestArticle(t *testing.T){
	meta := blog.Metadata{Title: "Learning is great - Doing is better", Categories : []string{"Thoughts"}, Date: "2021-11-04"}
	sut := blog.Article{Meta: meta}
	t.Run("write meta to io.Writer", func(t *testing.T) {
		var file bytes.Buffer
		sut.Write(&file)
		assert.Equal(t,meta.String(),file.String())
	})
}
