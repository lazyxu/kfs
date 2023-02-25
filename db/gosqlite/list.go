package gosqlite

import (
	"context"
	"github.com/lazyxu/kfs/db/dbBase"

	"github.com/lazyxu/kfs/dao"
)

func (db *DB) List(ctx context.Context, branchName string, splitPath []string) (dirItems []dao.DirItem, err error) {
	conn := db.getConn()
	defer db.putConn(conn)
	return dbBase.List(ctx, conn, branchName, splitPath)
}

func (db *DB) ListByHash(ctx context.Context, hash string) (dirItems []dao.DirItem, err error) {
	conn := db.getConn()
	defer db.putConn(conn)
	return dbBase.ListByHash(ctx, conn, hash)
}
