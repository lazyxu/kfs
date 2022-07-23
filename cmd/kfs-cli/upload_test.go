package main

import (
	"os"
	"strconv"
	"strings"
	"testing"

	sqlite "github.com/lazyxu/kfs/sqlite/noncgo"
)

func Test_upload_ls(t *testing.T) {
	// 1. reset
	{
		_, _ = exec(t, []string{"reset"})
	}
	// 2. upload file
	{
		stdout, stderr := exec(t, []string{"upload", "upload_test.go"})
		if stdout != "" {
			t.Errorf("expected \"\", actual \"%s\"", stdout)
		}
		if !strings.Contains(stderr, sqlite.ErrExpectedDir.Error()) {
			t.Errorf("expected \"%s\", actual \"%s\"", sqlite.ErrExpectedDir.Error(), stderr)
		}
	}
	// 3. upload dir
	{
		stdout, stderr := exec(t, []string{"upload", "."})
		if !strings.Contains(stdout, "commitId=2") {
			t.Errorf("expected \"%s\", actual \"%s\"", "commitId=2", stdout)
		}
		if stderr != "" {
			t.Errorf("expected \"\", actual \"%s\"", stderr)
		}
	}
	// 4. ls
	{
		stdout, stderr := exec(t, []string{"ls"})
		stdout = strings.Trim(stdout, "\n")
		items, err := os.ReadDir(".")
		if err != nil {
			t.Error(err)
		}
		expectedTotal := "total " + strconv.Itoa(len(items))
		if !strings.Contains(stdout, expectedTotal) {
			t.Errorf("expected \"%s\", actual \"%s\"", expectedTotal, stdout)
		}
		if stderr != "" {
			t.Errorf("expected \"\", actual \"%s\"", stderr)
		}
	}
}
