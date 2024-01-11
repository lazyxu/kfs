package cgosqlite

import (
	"context"
	"github.com/lazyxu/kfs/dao"
	"github.com/lazyxu/kfs/db/dbBase"
)

func (db *DB) ListDCIMDriver(ctx context.Context) (drivers []dao.DCIMDriver, err error) {
	conn := db.getConn()
	defer db.putConn(conn)
	return dbBase.ListDCIMDriver(ctx, conn)
}

func (db *DB) ListDCIMMediaType(ctx context.Context) (m map[string][]dao.Metadata, err error) {
	conn := db.getConn()
	defer db.putConn(conn)
	return dbBase.ListDCIMMediaType(ctx, conn)
}

func (db *DB) ListDCIMLocation(ctx context.Context) (list []dao.Metadata, err error) {
	conn := db.getConn()
	defer db.putConn(conn)
	return dbBase.ListDCIMLocation(ctx, conn)
}
