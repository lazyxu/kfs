package mysql

import (
	"context"
	"errors"

	"github.com/lazyxu/kfs/dao"
)

func (db *DB) ResetBranch(ctx context.Context, branchName string) error {
	conn := db.getConn()
	defer db.putConn(conn)
	dir, err := db.writeDir(ctx, conn, nil, nil)
	if err != nil {
		return err
	}
	commit := dao.NewCommit(dir, branchName, "")
	err = db.writeCommit(ctx, conn, &commit)
	if err != nil {
		return err
	}
	branch := dao.NewBranch(branchName, commit, dir)
	err = db.writeBranch(ctx, conn, branch)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) writeBranch(ctx context.Context, txOrDb TxOrDb, branch dao.Branch) error {
	_, err := txOrDb.ExecContext(ctx, `
	INSERT INTO _branch (
		name,
		description,
		commitId,
		size,
		count
	) VALUES (?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE 
		commitId=?,
		size=?,
		count=?;
	`, branch.Name, branch.Description, branch.CommitId, branch.Size, branch.Count,
		branch.CommitId, branch.Size, branch.Count)
	if err != nil {
		return err
	}
	return err
}

func (db *DB) NewBranch(ctx context.Context, branchName string) (exist bool, err error) {
	conn := db.getConn()
	defer db.putConn(conn)
	dir, err := db.writeDir(ctx, conn, nil, nil)
	if err != nil {
		return
	}
	commit := dao.NewCommit(dir, branchName, "")
	err = db.writeCommit(ctx, conn, &commit)
	if err != nil {
		return
	}
	branch := dao.NewBranch(branchName, commit, dir)
	err = db.writeBranch(ctx, conn, branch)
	if err != nil {
		return
	}
	return
}

func (db *DB) BranchInfo(ctx context.Context, branchName string) (branch dao.Branch, err error) {
	conn := db.getConn()
	defer db.putConn(conn)
	rows, err := conn.QueryContext(ctx, `
	SELECT * FROM _branch WHERE name=?;
	`, branchName)
	if err != nil {
		return
	}
	defer rows.Close()
	if !rows.Next() {
		return branch, errors.New("no such branch " + branchName)
	}
	err = rows.Scan(&branch.Name, &branch.Description, &branch.CommitId, &branch.Size, &branch.Count)
	if err != nil {
		return
	}
	return
}