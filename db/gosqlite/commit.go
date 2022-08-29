package gosqlite

import (
	"context"

	"github.com/lazyxu/kfs/dao"
)

func (db *DB) WriteCommit(ctx context.Context, commit *dao.Commit) error {
	conn := db.getConn()
	defer db.putConn(conn)
	return db.writeCommit(ctx, conn, commit)
}

func (db *DB) updateBranch(ctx context.Context, txOrDb TxOrDb, dir dao.Dir, branchName string, message string) error {
	commit := dao.NewCommit(dir, branchName, message)
	err := db.writeCommit(ctx, txOrDb, &commit)
	if err != nil {
		return err
	}
	branch := dao.NewBranch(branchName, commit, dir)
	err = db.insertBranch(ctx, txOrDb, branch)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) writeCommit(ctx context.Context, txOrDb TxOrDb, commit *dao.Commit) error {
	// TODO: if Hash not changed.
	res, err := txOrDb.ExecContext(ctx, `
	INSERT INTO _commit (createTime, Hash, lastId)
	VALUES (?, ?, ifnull((SELECT commitId FROM _branch WHERE _branch.name=?), 0));;
	`, commit.CreateTime(), commit.Hash, commit.BranchName())
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
