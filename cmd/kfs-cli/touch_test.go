package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTouch(t *testing.T) {
	filePath := "test"
	// 1. reset
	{
		_, _ = exec(t, []string{"reset"})
	}
	// 2. touch file
	{
		stdout, stderr := exec(t, []string{"touch", filePath})
		assert.Contains(t, stdout, "count=1")
		assert.Empty(t, stderr)
	}
	// 3. ls
	{
		stdout, stderr := exec(t, []string{"ls"})
		assert.Contains(t, stdout, "total 1")
		assert.Empty(t, stderr)
	}
	// 4. cat
	{
		stdout, stderr := exec(t, []string{"cat", filePath})
		assert.Empty(t, stdout)
		assert.Empty(t, stderr)
	}
}
