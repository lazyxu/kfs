package cgosqlite

import (
	"context"
	"github.com/lazyxu/kfs/dao"
	"github.com/lazyxu/kfs/db/dbBase"
)

func (db *DB) InsertDCIMMetadataTime(ctx context.Context, hash string, t int64) (exist bool, err error) {
	conn := db.getConn()
	defer db.putConn(conn)
	return dbBase.InsertDCIMMetadataTime(ctx, conn, db, hash, t)
}

func (db *DB) UpsertDCIMMetadataTime(ctx context.Context, hash string, t int64) error {
	conn := db.getConn()
	defer db.putConn(conn)
	return dbBase.UpsertDCIMMetadataTime(ctx, conn, hash, t)
}

func (db *DB) GetEarliestCrated(ctx context.Context, hash string) (t int64, err error) {
	conn := db.getConn()
	defer db.putConn(conn)
	return dbBase.GetEarliestCrated(ctx, conn, db, hash)
}

func (db *DB) ListMetadataTime(ctx context.Context) (list []dao.Metadata, err error) {
	conn := db.getConn()
	defer db.putConn(conn)
	return dbBase.ListMetadataTime(ctx, conn)
}
