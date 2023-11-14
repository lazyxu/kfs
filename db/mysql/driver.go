package mysql

import (
	"context"

	"github.com/lazyxu/kfs/dao"
	"github.com/lazyxu/kfs/db/dbBase"
)

func (db *DB) InsertDriver(ctx context.Context, driverName string, description string) (exist bool, err error) {
	return dbBase.InsertDriver(ctx, db.db, db, driverName, description)
}

func (db *DB) InsertDriverBaiduPhoto(ctx context.Context, driverName string, description string, accessToken string, refreshToken string) (exist bool, err error) {
	return dbBase.InsertDriverBaiduPhoto(ctx, db.db, db, driverName, description, accessToken, refreshToken)
}

func (db *DB) InsertDriverLocalFile(ctx context.Context, driverName string, description string, deviceId uint64, srcPath string, encoder string, concurrent int) (exist bool, err error) {
	return dbBase.InsertDriverLocalFile(ctx, db.db, db, driverName, description, deviceId, srcPath, encoder, concurrent)
}

func (db *DB) UpdateDriverSync(ctx context.Context, driverId uint64, sync bool, h int64, m int64) error {
	return dbBase.UpdateDriverSync(ctx, db.db, driverId, sync, h, m)
}

func (db *DB) DeleteDriver(ctx context.Context, driverId uint64) error {
	return dbBase.DeleteDriver(ctx, db.db, driverId)
}

func (db *DB) ListDriver(ctx context.Context) (drivers []dao.Driver, err error) {
	return dbBase.ListDriver(ctx, db.db)
}

func (db *DB) GetDriver(ctx context.Context, driverId uint64) (driver dao.Driver, err error) {
	return dbBase.GetDriver(ctx, db.db, driverId)
}

func (db *DB) GetDriverToken(ctx context.Context, driverId uint64) (driver dao.Driver, err error) {
	return dbBase.GetDriverToken(ctx, db.db, driverId)
}

func (db *DB) GetDriverSync(ctx context.Context, driverId uint64) (driver dao.Driver, err error) {
	return dbBase.GetDriverSync(ctx, db.db, driverId)
}

func (db *DB) ListCloudDriverSync(ctx context.Context) (drivers []dao.Driver, err error) {
	return dbBase.ListCloudDriverSync(ctx, db.db)
}

func (db *DB) GetDriverLocalFile(ctx context.Context, driverId uint64) (driver dao.Driver, err error) {
	return dbBase.GetDriverLocalFile(ctx, db.db, driverId)
}

func (db *DB) GetDriverFileSize(ctx context.Context, driverId uint64) (n uint64, err error) {
	return dbBase.GetDriverFileSize(ctx, db.db, driverId)
}

func (db *DB) GetDriverFileCount(ctx context.Context, driverId uint64) (n uint64, err error) {
	return dbBase.GetDriverFileCount(ctx, db.db, driverId)
}

func (db *DB) GetDriverDirCount(ctx context.Context, driverId uint64) (n uint64, err error) {
	return dbBase.GetDriverDirCount(ctx, db.db, driverId)
}
