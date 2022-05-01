package blog

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)


func TestBookFilePath(t *testing.T) {
	sut := NewBookWithBaseDir(Metadata{Title: "Book title"},"/b")
	assert.Equal(t, filepath.Join("/b",bookDir,"Book title.md"), sut.Path())
}
