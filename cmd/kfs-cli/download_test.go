package main

import (
	"io/ioutil"
	"os"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

var downloadStdoutRegex = regexp.MustCompile("Saving to '(.+)'")

func TestDownloadFile(t *testing.T) {
	fileName := "download_test.go"
	// 1. reset
	{
		_, _ = exec(t, []string{"reset"})
	}
	// 2. upload file
	{
		stdout, stderr := exec(t, []string{"upload", fileName, "-p", fileName})
		assert.Contains(t, stdout, "commitId=2")
		assert.Empty(t, stderr)
	}
	// 3. download file
	{
		tempFilePath := os.TempDir() + fileName
		_ = os.Remove(tempFilePath)
		stdout, stderr := exec(t, []string{"download", fileName, "-p", tempFilePath})
		assert.Regexp(t, downloadStdoutRegex, stdout)
		filePath := downloadStdoutRegex.FindStringSubmatch(stdout)
		assert.Equal(t, 2, len(filePath))
		assert.Empty(t, stderr)
		expected, err := ioutil.ReadFile(fileName)
		assert.Nil(t, err)
		actual, err := ioutil.ReadFile(filePath[1])
		assert.Nil(t, err)
		assert.Equal(t, expected, actual)
	}
}

func TestDownloadDir(t *testing.T) {
	// 1. reset
	{
		_, _ = exec(t, []string{"reset"})
	}
	// 2. upload dir
	{
		stdout, stderr := exec(t, []string{"upload", "."})
		assert.Contains(t, stdout, "commitId=2")
		assert.Empty(t, stderr)
	}
	// 3. download dir
	{
		fileName := "kfs-test"
		tempFilePath := os.TempDir() + fileName
		_ = os.Remove(tempFilePath)
		stdout, stderr := exec(t, []string{"download", "-p", tempFilePath})
		assert.Empty(t, stderr)
		assert.Regexp(t, downloadStdoutRegex, stdout)
		filePath := downloadStdoutRegex.FindStringSubmatch(stdout)
		assert.Equal(t, 2, len(filePath))
		tempItems, err := os.ReadDir(filePath[1])
		assert.Nil(t, err)
		items, err := os.ReadDir(".")
		assert.Nil(t, err)
		assert.Equal(t, len(items), len(tempItems))
	}
}
