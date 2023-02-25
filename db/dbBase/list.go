package dbBase

import (
	"context"
	"database/sql"
	"github.com/lazyxu/kfs/dao"
)

func List(ctx context.Context, conn *sql.DB, branchName string, splitPath []string) (dirItems []dao.DirItem, err error) {
	tx, err := conn.Begin()
	if err != nil {
		return
	}
	defer func() {
		err = CommitAndRollback(tx, err)
	}()
	hash, err := getBranchCommitHash(ctx, tx, branchName)
	if err != nil {
		return
	}
	for i := range splitPath {
		hash, err = getDirItemHash(ctx, tx, hash, splitPath, i)
		if err != nil {
			return
		}
	}
	dirItems, err = getDirItems(ctx, tx, hash)
	if err != nil {
		return
	}
	return
}

func ListByHash(ctx context.Context, conn *sql.DB, hash string) (dirItems []dao.DirItem, err error) {
	tx, err := conn.Begin()
	if err != nil {
		return
	}
	defer func() {
		err = CommitAndRollback(tx, err)
	}()
	dirItems, err = getDirItems(ctx, tx, hash)
	if err != nil {
		return
	}
	return
}
