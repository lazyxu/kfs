package gosqlite

import (
	"context"
	"github.com/lazyxu/kfs/dao"
	"github.com/lazyxu/kfs/db/dbBase"
)

func (db *DB) MaxBatchSize() int {
	return 32766
}

func (db *DB) WriteDir(ctx context.Context, dirItems []dao.DirItem) (dir dao.Dir, err error) {
	conn := db.getConn()
	defer db.putConn(conn)
	return dbBase.WriteDir(ctx, conn, db, dirItems)
}

func (db *DB) RemoveDirItem(ctx context.Context, branchName string, splitPath []string) (commit dao.Commit, branch dao.Branch, err error) {
	conn := db.getConn()
	defer db.putConn(conn)
	return dbBase.RemoveDirItem(ctx, conn, db, branchName, splitPath)
}
