package gosqlite

import (
	"context"
	"os"

	"github.com/lazyxu/kfs/dao"
)

func (db *DB) Open(ctx context.Context, branchName string, splitPath []string) (hash string, mode os.FileMode, dirItems []dao.DirItem, err error) {
	conn := db.getConn()
	defer db.putConn(conn)
	tx, err := conn.Begin()
	if err != nil {
		return
	}
	defer func() {
		err = commitAndRollback(tx, err)
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

func (db *DB) Open2(ctx context.Context, branchName string, splitPath []string) (dirItem dao.DirItem, dirItems []dao.DirItem, err error) {
	conn := db.getConn()
	defer db.putConn(conn)
	tx, err := conn.Begin()
	if err != nil {
		return
	}
	defer func() {
		err = commitAndRollback(tx, err)
	}()
	hash, err := db.getBranchCommitHash(ctx, tx, branchName)
	if err != nil {
		return
	}
	if len(splitPath) != 0 {
		for i := range splitPath {
			dirItem, err = db.getDirItem(ctx, tx, hash, splitPath, i)
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
	dirItems, err = db.getDirItems(ctx, tx, hash)
	if err != nil {
		return
	}
	return
}
