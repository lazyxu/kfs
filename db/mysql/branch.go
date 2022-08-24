package mysql

import (
	"context"
	"errors"

	"github.com/lazyxu/kfs/dao"
)

func (db *DB) WriteBranch(ctx context.Context, branch dao.Branch) error {
	conn := db.getConn()
	defer db.putConn(conn)
	return db.writeBranch(ctx, conn, branch)
}

func (db *DB) writeBranch(ctx context.Context, txOrDb TxOrDb, branch dao.Branch) error {
	_, err := txOrDb.ExecContext(ctx, `
	REPLACE INTO branch VALUES (?, ?, ?, ?, ?);
	`, branch.Name, branch.Description, branch.CommitId, branch.Size, branch.Count)
	return err
}

func (db *DB) insertBranch(ctx context.Context, txOrDb TxOrDb, branch dao.Branch) error {
	_, err := txOrDb.ExecContext(ctx, `
	INSERT INTO branch (
		name,
		commitId,
		size,
		count
	) VALUES (?, ?, ?, ?);
	`, branch.Name, branch.CommitId, branch.Size, branch.Count)
	return err
}

func (db *DB) NewBranch(ctx context.Context, branchName string) (exist bool, err error) {
	conn := db.getConn()
	defer db.putConn(conn)
	tx, err := conn.Begin()
	if err != nil {
		return
	}
	defer func() {
		err = commitAndRollback(tx, err)
	}()
	dir, err := db.writeDir(ctx, tx, nil, nil)
	if err != nil {
		return
	}
	commit := dao.NewCommit(dir, branchName, "")
	err = db.writeCommit(ctx, tx, &commit)
	if err != nil {
		return
	}
	branch := dao.NewBranch(branchName, commit, dir)
	err = db.insertBranch(ctx, tx, branch)
	if isUniqueConstraintError(err) {
		exist = true
		err = nil
	}
	return
}

func (db *DB) BranchInfo(ctx context.Context, branchName string) (branch dao.Branch, err error) {
	conn := db.getConn()
	defer db.putConn(conn)
	rows, err := conn.QueryContext(ctx, `
	SELECT * FROM branch WHERE name=?;
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
