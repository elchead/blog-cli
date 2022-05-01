package blog

import (
	"bytes"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)


func TestArticleFilePath(t *testing.T) {
	sut := NewArticle(Metadata{Title: "Article title"})
	assert.Equal(t, filepath.Join(obsidianVault,articleDir, "Article title.md"), sut.Path())
}


func TestArticle(t *testing.T){
	meta := Metadata{Title: "Learning is great - Doing is better", Categories : []string{"Thoughts"}, Date: "2021-11-04"}
	sut := Article{Meta: meta}
	t.Run("write meta to io.Writer", func(t *testing.T) {
		var file bytes.Buffer
		sut.Write(&file)
		assert.Equal(t,meta.String(),file.String())
	})
}
