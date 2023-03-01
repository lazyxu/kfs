package dbBase

import (
	"context"
	"database/sql"
	"errors"
	"github.com/lazyxu/kfs/dao"
)

func ResetBranch(ctx context.Context, conn *sql.DB, db DbImpl, branchName string) (err error) {
	tx, err := conn.Begin()
	if err != nil {
		return err
	}
	defer func() {
		err = CommitAndRollback(tx, err)
	}()
	dir, err := InsertDirWithTx(ctx, tx, db, nil, nil)
	if err != nil {
		return err
	}
	commit := dao.NewCommit(dir, branchName, "")
	err = db.InsertCommitWithTxOrDb(ctx, conn, &commit)
	if err != nil {
		return err
	}
	branch := dao.NewBranch(branchName, commit, dir)
	err = db.UpsertBranchWithTxOrDb(ctx, conn, branch)
	if err != nil {
		return err
	}
	return nil
}

func NewBranch(ctx context.Context, conn *sql.DB, db DbImpl, branchName string) (exist bool, err error) {
	tx, err := conn.Begin()
	if err != nil {
		return
	}
	defer func() {
		err = CommitAndRollback(tx, err)
	}()
	dir, err := InsertDirWithTx(ctx, tx, db, nil, nil)
	if err != nil {
		return
	}
	commit := dao.NewCommit(dir, branchName, "")
	err = db.InsertCommitWithTxOrDb(ctx, tx, &commit)
	if err != nil {
		return
	}
	branch := dao.NewBranch(branchName, commit, dir)
	err = InsertBranchWithTxOrDb(ctx, tx, branch)
	if db.IsUniqueConstraintError(err) {
		exist = true
		err = nil
	}
	return
}

func InsertBranchWithTxOrDb(ctx context.Context, txOrDb TxOrDb, branch dao.Branch) error {
	_, err := txOrDb.ExecContext(ctx, `
	INSERT INTO _branch (
		name,
		description,
		commitId,
		size,
		count
	) VALUES (?, ?, ?, ?, ?)`, branch.Name, branch.Description, branch.CommitId, branch.Size, branch.Count)
	if err != nil {
		return err
	}
	return err
}

func DeleteBranch(ctx context.Context, conn *sql.DB, name string) error {
	_, err := conn.ExecContext(ctx, `
	DELETE FROM _branch WHERE name = ?`, name)
	if err != nil {
		return err
	}
	return err
}

func UpsertBranchWithTxOrDb(ctx context.Context, txOrDb TxOrDb, branch dao.Branch) error {
	_, err := txOrDb.ExecContext(ctx, `
	INSERT INTO _branch (
		name,
		description,
		commitId,
		size,
		count
	) VALUES (?, ?, ?, ?, ?) ON CONFLICT(name) DO UPDATE SET
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

func UpsertBranchWithTxOrDbMysql(ctx context.Context, txOrDb TxOrDb, branch dao.Branch) error {
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

func BranchInfo(ctx context.Context, txOrDb TxOrDb, branchName string) (branch dao.Branch, err error) {
	rows, err := txOrDb.QueryContext(ctx, `
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

func BranchList(ctx context.Context, txOrDb TxOrDb) (branches []dao.IBranch, err error) {
	rows, err := txOrDb.QueryContext(ctx, `
	SELECT * FROM _branch;
	`)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var branch dao.Branch
		err = rows.Scan(&branch.Name, &branch.Description, &branch.CommitId, &branch.Size, &branch.Count)
		if err != nil {
			return
		}
		branches = append(branches, branch)
	}
	return
}
