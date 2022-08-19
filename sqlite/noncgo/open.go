package noncgo

import (
	"context"
	"os"
)

func (db *DB) Open(ctx context.Context, branchName string, splitPath []string) (hash string, mode os.FileMode, dirItems []DirItem, err error) {
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
	hash, err = db.getBranchCommitHash(ctx, tx, branchName)
	if err != nil {
		return
	}
	var m uint64
	for i := range splitPath {
		hash, m, err = db.getDirItemHashMode(ctx, tx, hash, splitPath, i)
		if err != nil {
			return
		}
	}
	mode = os.FileMode(m)
	if mode.IsDir() {
		dirItems, err = db.getDirItems(ctx, tx, hash)
		if err != nil {
			return
		}
	}
	return
}
