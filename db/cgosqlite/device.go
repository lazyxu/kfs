package cgosqlite

import (
	"context"
	"github.com/lazyxu/kfs/dao"
	"github.com/lazyxu/kfs/db/dbBase"
)

func (db *DB) InsertDevice(ctx context.Context, name string, os string) (int64, error) {
	conn := db.getConn()
	defer db.putConn(conn)
	return dbBase.InsertDevice(ctx, conn, name, os)
}

func (db *DB) DeleteDevice(ctx context.Context, deviceId uint64) error {
	conn := db.getConn()
	defer db.putConn(conn)
	return dbBase.DeleteDevice(ctx, conn, deviceId)
}

func (db *DB) ListDevice(ctx context.Context) (devices []dao.Device, err error) {
	conn := db.getConn()
	defer db.putConn(conn)
	return dbBase.ListDevice(ctx, conn)
}
