package main

import (
	"fmt"
	"testing"

	"github.com/skratchdot/open-golang/open"
	"github.com/stretchr/testify/assert"
)

func Test(t *testing.T){
	err := open.Run("obsidian://open?file=Money.md")
	fmt.Println(err)
	assert.NoError(t, err)
}
