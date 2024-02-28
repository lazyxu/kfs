package gosqlite

import (
	"context"
	"github.com/lazyxu/kfs/dao"
	"github.com/lazyxu/kfs/db/dbBase"
)

func (db *DB) InsertDevice(ctx context.Context, id string, name string, os string, userAgent string, hostname string) error {
	conn := db.getConn()
	defer db.putConn(conn)
	return dbBase.InsertDevice(ctx, conn, id, name, os, userAgent, hostname)
}

func (db *DB) DeleteDevice(ctx context.Context, deviceId string) error {
	conn := db.getConn()
	defer db.putConn(conn)
	return dbBase.DeleteDevice(ctx, conn, deviceId)
}

func (db *DB) ListDevice(ctx context.Context) (devices []dao.Device, err error) {
	conn := db.getConn()
	defer db.putConn(conn)
	return dbBase.ListDevice(ctx, conn)
}
