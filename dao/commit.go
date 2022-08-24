package dao

import (
	"time"
)

type Commit struct {
	Id         uint64
	createTime uint64
	Hash       string
	lastId     uint64
	branchName string
}

func NewCommit(dir Dir, branchName string, message string) Commit {
	return Commit{0, uint64(time.Now().UnixNano()), dir.hash, 0, branchName}
}

func (c *Commit) CreateTime() uint64 {
	return c.createTime
}

func (c *Commit) BranchName() string {
	return c.branchName
}
