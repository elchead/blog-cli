package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAskToPublish(t *testing.T) {
	t.Run("input yes",func(t *testing.T) {
		reader := strings.NewReader("y!\n")
		assert.True(t,okToPublish(reader))
	})
	t.Run("input no",func(t *testing.T) {
		reader := strings.NewReader("n!\n")
		assert.False(t,okToPublish(reader))
	})
	t.Run("input wrong yes",func(t *testing.T) {
		reader := strings.NewReader("y\n")
		assert.False(t,okToPublish(reader))
	})
}

