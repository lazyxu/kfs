package noncgo

import (
	"context"
	"errors"

	"modernc.org/sqlite"
)

type Branch struct {
	Name        string
	Description string
	CommitId    uint64
	Size        uint64
	Count       uint64
}

func NewBranch(name string, description string, commit Commit, dir Dir) Branch {
	return Branch{name, description, commit.Id, dir.size, dir.count}
}

func (db *DB) WriteBranch(ctx context.Context, branch Branch) error {
	return db.writeBranch(ctx, db._db, branch)
}

func (db *DB) writeBranch(ctx context.Context, txOrDb TxOrDb, branch Branch) error {
	_, err := txOrDb.ExecContext(ctx, `
	REPLACE INTO branch VALUES (?, ?, ?, ?, ?);
	`, branch.Name, branch.Description, branch.CommitId, branch.Size, branch.Count)
	return err
}

func (db *DB) insertBranch(ctx context.Context, txOrDb TxOrDb, branch Branch) error {
	_, err := txOrDb.ExecContext(ctx, `
	INSERT INTO branch VALUES (?, ?, ?, ?, ?);
	`, branch.Name, branch.Description, branch.CommitId, branch.Size, branch.Count)
	return err
}

func (db *DB) NewBranch(ctx context.Context, branchName string, description string) (exist bool, err error) {
	tx, err := db._db.Begin()
	if err != nil {
		return
	}
	defer func() {
		if err != nil {
			return
		}
		err = tx.Commit()
		if err == nil {
			return
		}
		e, ok := err.(*sqlite.Error)
		if ok && e.Code() == 5 {
			return
		}
		err1 := tx.Rollback()
		if err1 != nil {
			panic(err1) // should not happen
		}
		return
	}()
	dir, err := db.writeDir(ctx, tx, nil)
	if err != nil {
		return
	}
	commit := NewCommit(dir, branchName)
	err = db.writeCommit(ctx, tx, &commit)
	if err != nil {
		return
	}
	branch := NewBranch(branchName, description, commit, dir)
	err = db.insertBranch(ctx, tx, branch)
	if isUniqueConstraintError(err) {
		exist = true
		err = nil
	}
	return
}

func (db *DB) BranchInfo(ctx context.Context, branchName string) (branch Branch, err error) {
	rows, err := db._db.QueryContext(ctx, `
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
