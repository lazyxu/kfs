package noncgo

import (
	"context"

	"github.com/lazyxu/kfs/dao"
)

func (db *DB) List(ctx context.Context, branchName string, splitPath []string) (dirItems []dao.DirItem, err error) {
	conn := db.getConn()
	defer db.putConn(conn)
	tx, err := conn.Begin()
	if err != nil {
		return
	}
	defer func() {
		if err == nil {
			err = tx.Commit()
			if err != nil {
				err1 := tx.Rollback()
				if err1 != nil {
					panic(err1) // should not happen
				}
				return
			}
		}
	}()
	hash, err := db.getBranchCommitHash(ctx, tx, branchName)
	if err != nil {
		return
	}
	for i := range splitPath {
		hash, err = db.getDirItemHash(ctx, tx, hash, splitPath, i)
		if err != nil {
			return
		}
	}
	dirItems, err = db.getDirItems(ctx, tx, hash)
	if err != nil {
		return
	}
	return
}

func (db *DB) ListByHash(ctx context.Context, hash string) (dirItems []dao.DirItem, err error) {
	conn := db.getConn()
	defer db.putConn(conn)
	tx, err := conn.Begin()
	if err != nil {
		return
	}
	defer func() {
		if err == nil {
			err = tx.Commit()
			if err != nil {
				err1 := tx.Rollback()
				if err1 != nil {
					panic(err1) // should not happen
				}
				return
			}
		}
	}()
	dirItems, err = db.getDirItems(ctx, tx, hash)
	if err != nil {
		return
	}
	return
}
