package main

import (
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBranch(t *testing.T) {
	// 1. reset
	{
		_, _ = exec(t, []string{"reset"})
	}
	// 2. branch info
	{
		stdout, stderr := exec(t, []string{"branch", "info", "master"})
		assert.Contains(t, stdout, "size: 0\ncount: 0\n")
		assert.Empty(t, stderr)
	}
	branchName := strconv.Itoa(int(time.Now().UnixNano()))
	// 3. checkout branch
	{
		stdout, stderr := exec(t, []string{"branch", "checkout", branchName})
		assert.Contains(t, stdout, "switch to branch '"+branchName+"'")
		assert.Empty(t, stderr)
	}
	// 4. branch info
	{
		stdout, stderr := exec(t, []string{"branch", "info", branchName})
		assert.Contains(t, stdout, "size: 0\ncount: 0\n")
		assert.Empty(t, stderr)
	}
}
