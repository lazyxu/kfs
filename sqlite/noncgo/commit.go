package noncgo

import (
	"context"
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

func (db *DB) WriteCommit(ctx context.Context, commit *Commit) error {
	conn := db.getConn()
	defer db.putConn(conn)
	return db.writeCommit(ctx, conn, commit)
}

func (db *DB) writeCommit(ctx context.Context, txOrDb TxOrDb, commit *Commit) error {
	// TODO: if Hash not changed.
	res, err := txOrDb.ExecContext(ctx, `
	INSERT INTO [commit] (createTime, Hash, lastId)
	VALUES (?, ?, ifnull((SELECT commitId FROM branch WHERE branch.name=?), 0));;
	`, commit.createTime, commit.Hash, commit.branchName)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	commit.Id = uint64(id)
	return err
}
