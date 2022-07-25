package main

import (
	"os"
	"strconv"
	"testing"

	sqlite "github.com/lazyxu/kfs/sqlite/noncgo"

	"github.com/stretchr/testify/assert"
)

func TestUpload(t *testing.T) {
	// 1. reset
	{
		_, _ = exec(t, []string{"reset"})
	}
	// 2. upload file
	{
		stdout, stderr := exec(t, []string{"upload", "upload_test.go"})
		assert.Empty(t, stdout)
		assert.Contains(t, stderr, sqlite.ErrExpectedDir.Error())
	}
	// 3. upload dir
	{
		stdout, stderr := exec(t, []string{"upload", "."})
		assert.Contains(t, stdout, "commitId=2")
		assert.Empty(t, stderr)
	}
	// 4. ls
	{
		stdout, stderr := exec(t, []string{"ls"})
		items, err := os.ReadDir(".")
		assert.Nil(t, err)
		assert.Contains(t, stdout, "total "+strconv.Itoa(len(items)))
		assert.Empty(t, stderr)
	}
}
