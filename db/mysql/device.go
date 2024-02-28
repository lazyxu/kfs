package mysql

import (
	"context"
	"github.com/lazyxu/kfs/dao"
	"github.com/lazyxu/kfs/db/dbBase"
)

func (db *DB) InsertDevice(ctx context.Context, id string, name string, os string, userAgent string, hostname string) error {
	return dbBase.InsertDevice(ctx, db.db, id, name, os, userAgent, hostname)
}

func (db *DB) DeleteDevice(ctx context.Context, deviceId string) error {
	return dbBase.DeleteDevice(ctx, db.db, deviceId)
}

func (db *DB) ListDevice(ctx context.Context) (devices []dao.Device, err error) {
	return dbBase.ListDevice(ctx, db.db)
}
