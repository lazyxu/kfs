package mysql

import (
	"context"
	"github.com/lazyxu/kfs/dao"
	"github.com/lazyxu/kfs/db/dbBase"
)

func (db *DB) InsertDevice(ctx context.Context, name string, os string) (int64, error) {
	return dbBase.InsertDevice(ctx, db.db, name, os)
}

func (db *DB) DeleteDevice(ctx context.Context, deviceId uint64) error {
	return dbBase.DeleteDevice(ctx, db.db, deviceId)
}

func (db *DB) ListDevice(ctx context.Context) (devices []dao.Device, err error) {
	return dbBase.ListDevice(ctx, db.db)
}
