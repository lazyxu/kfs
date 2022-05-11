package noncgo

import (
	"context"
	"time"
)

type Commit struct {
	id         int64
	createTime uint64
	hash       string
	lastId     uint64
	branchName string
}

func NewCommit(dir Dir, branchName string) Commit {
	return Commit{0, uint64(time.Now().UnixNano()), dir.hash, 0, branchName}
}

func (db *SqliteNonCgoDB) WriteCommit(ctx context.Context, commit *Commit) error {
	res, err := db._db.ExecContext(ctx, `
	INSERT INTO [commit] (createTime, hash, lastId)
	SELECT ?, ?, commitId FROM branch WHERE branch.name=?;
	`, commit.createTime, commit.hash, commit.branchName)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	commit.id = id
	return err
}
