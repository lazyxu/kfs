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

func (db *DB) WriteCommit(ctx context.Context, commit *Commit) error {
	return db.writeCommit(ctx, db._db, commit)
}

func (db *DB) writeCommit(ctx context.Context, txOrDb TxOrDb, commit *Commit) error {
	res, err := txOrDb.ExecContext(ctx, `
	INSERT INTO [commit] (createTime, hash, lastId)
	VALUES (?, ?, ifnull((SELECT commitId FROM branch WHERE branch.name=?), 0));;
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
