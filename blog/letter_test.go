package blog

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)



func TestLetterFilePath(t *testing.T) {
	sut := NewLetter(Metadata{Title: "Letter title"})
	assert.Equal(t, filepath.Join(letterDir,"Letter title.md"), sut.Path())
}
