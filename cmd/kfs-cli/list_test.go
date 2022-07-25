package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestList(t *testing.T) {
	// 1. reset
	{
		_, _ = exec(t, []string{"reset"})
	}
	// 2. ls
	{
		stdout, stderr := exec(t, []string{"ls"})
		assert.Contains(t, stdout, "total 0")
		assert.Empty(t, stderr)
	}
}
