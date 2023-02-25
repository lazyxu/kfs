package dbBase

import (
	"context"
	"database/sql"
	"github.com/lazyxu/kfs/dao"
	"os"
)

func Open(ctx context.Context, conn *sql.DB, branchName string, splitPath []string) (hash string, mode os.FileMode, dirItems []dao.DirItem, err error) {
	tx, err := conn.Begin()
	if err != nil {
		return
	}
	defer func() {
		err = CommitAndRollback(tx, err)
	}()
	hash, err = getBranchCommitHash(ctx, tx, branchName)
	if err != nil {
		return
	}
	var m uint64
	for i := range splitPath {
		hash, m, err = getDirItemHashMode(ctx, tx, hash, splitPath, i)
		if err != nil {
			return
		}
	}
	mode = os.FileMode(m)
	if mode.IsDir() {
		dirItems, err = getDirItems(ctx, tx, hash)
		if err != nil {
			return
		}
	}
	return
}

func Open2(ctx context.Context, conn *sql.DB, branchName string, splitPath []string) (dirItem dao.DirItem, dirItems []dao.DirItem, err error) {
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
	if len(splitPath) != 0 {
		for i := range splitPath {
			dirItem, err = getDirItem(ctx, tx, hash, splitPath, i)
			if err != nil {
				return
			}
			hash = dirItem.Hash
		}
		mode := os.FileMode(dirItem.Mode)
		if mode.IsRegular() {
			return
		}
	}
	dirItem.Mode = uint64(os.ModeDir | os.ModePerm)
	dirItems, err = getDirItems(ctx, tx, hash)
	if err != nil {
		return
	}
	return
}
