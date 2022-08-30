package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCatFile(t *testing.T) {
	fileName := "cat_test.go"
	// 1. reset
	{
		_, _ = exec(t, []string{"reset"})
	}
	// 2. upload file
	{
		stdout, stderr := exec(t, []string{"upload", fileName, "-p", fileName})
		assert.Contains(t, stdout, "count=1")
		assert.Empty(t, stderr)
	}
	// 3. cat file
	{
		stdout, stderr := exec(t, []string{"cat", fileName})
		assert.Empty(t, stderr)
		expected, err := os.ReadFile(fileName)
		assert.Nil(t, err)
		assert.Equal(t, string(expected), stdout)
	}
}
