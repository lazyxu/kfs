package mysql

import (
	"context"
	"github.com/lazyxu/kfs/dao"
	"github.com/lazyxu/kfs/db/dbBase"
)

func (db *DB) InsertDriver(ctx context.Context, driverName string, description string, typ string, accessToken string, refreshToken string) (exist bool, err error) {
	return dbBase.InsertDriver(ctx, db.db, db, driverName, description, typ, accessToken, refreshToken)
}

func (db *DB) UpdateDriverSync(ctx context.Context, driverName string, sync bool, h int64, m int64, s int64) error {
	return dbBase.UpdateDriverSync(ctx, db.db, driverName, sync, h, m, s)
}

func (db *DB) DeleteDriver(ctx context.Context, driverName string) error {
	return dbBase.DeleteDriver(ctx, db.db, driverName)
}

func (db *DB) ListDriver(ctx context.Context) (drivers []dao.Driver, err error) {
	return dbBase.ListDriver(ctx, db.db)
}

func (db *DB) GetDriverToken(ctx context.Context, driverName string) (driver dao.Driver, err error) {
	return dbBase.GetDriverToken(ctx, db.db, driverName)
}

func (db *DB) GetDriverSync(ctx context.Context, driverName string) (driver dao.Driver, err error) {
	return dbBase.GetDriverSync(ctx, db.db, driverName)
}

func (db *DB) GetDriverFileSize(ctx context.Context, driverName string) (n uint64, err error) {
	return dbBase.GetDriverFileSize(ctx, db.db, driverName)
}

func (db *DB) GetDriverFileCount(ctx context.Context, driverName string) (n uint64, err error) {
	return dbBase.GetDriverFileCount(ctx, db.db, driverName)
}

func (db *DB) GetDriverDirCount(ctx context.Context, driverName string) (n uint64, err error) {
	return dbBase.GetDriverDirCount(ctx, db.db, driverName)
}
