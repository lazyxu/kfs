package cgosqlite

import (
	"context"
	"github.com/lazyxu/kfs/dao"
	"github.com/lazyxu/kfs/db/dbBase"
)

func (db *DB) InsertFileType(ctx context.Context, hash string, t dao.FileType) (exist bool, err error) {
	conn := db.getConn()
	defer db.putConn(conn)
	return dbBase.InsertFileType(ctx, conn, db, hash, t)
}

func (db *DB) ListExpectFileType(ctx context.Context) (hashList []string, err error) {
	conn := db.getConn()
	defer db.putConn(conn)
	return dbBase.ListExpectFileType(ctx, conn)
}

func (db *DB) ListFileType(ctx context.Context) (fileTypeMap map[string]dao.FileType, err error) {
	conn := db.getConn()
	defer db.putConn(conn)
	return dbBase.ListFileType(ctx, conn)
}
