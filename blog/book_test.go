package blog

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)


func TestBookFilePath(t *testing.T) {
	sut := NewBook(Metadata{Title: "Book title"})
	assert.Equal(t, filepath.Join(bookDir,"Book title.md"), sut.Path())
}
