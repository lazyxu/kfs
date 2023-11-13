package cgosqlite

import (
	"context"
	"github.com/lazyxu/kfs/dao"
	"github.com/lazyxu/kfs/db/dbBase"
)

func (db *DB) InsertDriver(ctx context.Context, driverName string, description string) (exist bool, err error) {
	conn := db.getConn()
	defer db.putConn(conn)
	return dbBase.InsertDriver(ctx, conn, db, driverName, description)
}

func (db *DB) InsertDriverBaiduPhoto(ctx context.Context, driverName string, description string, accessToken string, refreshToken string) (exist bool, err error) {
	conn := db.getConn()
	defer db.putConn(conn)
	return dbBase.InsertDriverBaiduPhoto(ctx, conn, db, driverName, description, accessToken, refreshToken)
}

func (db *DB) InsertDriverLocalFile(ctx context.Context, driverName string, description string, deviceId uint64, srcPath string, encoder string, concurrent int) (exist bool, err error) {
	conn := db.getConn()
	defer db.putConn(conn)
	return dbBase.InsertDriverLocalFile(ctx, conn, db, driverName, description, deviceId, srcPath, encoder, concurrent)
}

func (db *DB) UpdateDriverSync(ctx context.Context, driverId uint64, sync bool, h int64, m int64, s int64) error {
	conn := db.getConn()
	defer db.putConn(conn)
	return dbBase.UpdateDriverSync(ctx, conn, driverId, sync, h, m, s)
}

func (db *DB) DeleteDriver(ctx context.Context, driverId uint64) error {
	conn := db.getConn()
	defer db.putConn(conn)
	return dbBase.DeleteDriver(ctx, conn, driverId)
}

func (db *DB) ListDriver(ctx context.Context) (drivers []dao.Driver, err error) {
	conn := db.getConn()
	defer db.putConn(conn)
	return dbBase.ListDriver(ctx, conn)
}

func (db *DB) GetDriverToken(ctx context.Context, driverId uint64) (driver dao.Driver, err error) {
	conn := db.getConn()
	defer db.putConn(conn)
	return dbBase.GetDriverToken(ctx, conn, driverId)
}

func (db *DB) GetDriverSync(ctx context.Context, driverId uint64) (driver dao.Driver, err error) {
	conn := db.getConn()
	defer db.putConn(conn)
	return dbBase.GetDriverSync(ctx, conn, driverId)
}

func (db *DB) GetDriverLocalFile(ctx context.Context, driverId uint64) (driver dao.Driver, err error) {
	conn := db.getConn()
	defer db.putConn(conn)
	return dbBase.GetDriverLocalFile(ctx, conn, driverId)
}

func (db *DB) GetDriverFileSize(ctx context.Context, driverId uint64) (n uint64, err error) {
	conn := db.getConn()
	defer db.putConn(conn)
	return dbBase.GetDriverFileSize(ctx, conn, driverId)
}

func (db *DB) GetDriverFileCount(ctx context.Context, driverId uint64) (n uint64, err error) {
	conn := db.getConn()
	defer db.putConn(conn)
	return dbBase.GetDriverFileCount(ctx, conn, driverId)
}

func (db *DB) GetDriverDirCount(ctx context.Context, driverId uint64) (n uint64, err error) {
	conn := db.getConn()
	defer db.putConn(conn)
	return dbBase.GetDriverDirCount(ctx, conn, driverId)
}
