package dao

import (
	"time"
)

type Commit struct {
	Id         uint64
	createTime time.Time
	Hash       string
	lastId     uint64
	branchName string
}

func NewCommit(dir Dir, branchName string, message string) Commit {
	return Commit{0, time.Now(), dir.hash, 0, branchName}
}

func (c *Commit) CreateTime() time.Time {
	return c.createTime
}

func (c *Commit) BranchName() string {
	return c.branchName
}
