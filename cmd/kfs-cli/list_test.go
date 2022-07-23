package main

import (
	"strings"
	"testing"
)

func Test_init_ls(t *testing.T) {
	// 1. reset
	{
		_, _ = exec(t, []string{"reset"})
	}
	// 2. ls
	{
		stdout, stderr := exec(t, []string{"ls"})
		stdout = strings.Trim(stdout, "\n")
		if stdout != "total 0" {
			t.Errorf("expected \"total 0\", actual \"%s\"", stdout)
		}
		if stderr != "" {
			t.Errorf("expected \"\", actual \"%s\"", stderr)
		}
	}
}
