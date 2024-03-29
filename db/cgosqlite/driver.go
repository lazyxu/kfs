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

func (db *DB) InsertDriverLocalFile(ctx context.Context, driverName string, description string, deviceId string, srcPath string, ignores string, encoder string) (exist bool, err error) {
	conn := db.getConn()
	defer db.putConn(conn)
	return dbBase.InsertDriverLocalFile(ctx, conn, db, driverName, description, deviceId, srcPath, ignores, encoder)
}

func (db *DB) UpdateDriverSync(ctx context.Context, driverId uint64, sync bool, h int64, m int64) error {
	conn := db.getConn()
	defer db.putConn(conn)
	return dbBase.UpdateDriverSync(ctx, conn, driverId, sync, h, m)
}

func (db *DB) UpdateDriverLocalFile(ctx context.Context, driverId uint64, srcPath, ignores, encoder string) error {
	conn := db.getConn()
	defer db.putConn(conn)
	return dbBase.UpdateDriverLocalFile(ctx, conn, driverId, srcPath, ignores, encoder)
}

func (db *DB) ResetDriver(ctx context.Context, driverId uint64) error {
	conn := db.getConn()
	defer db.putConn(conn)
	return dbBase.ResetDriver(ctx, conn, driverId)
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

func (db *DB) GetDriver(ctx context.Context, driverId uint64) (driver dao.Driver, err error) {
	conn := db.getConn()
	defer db.putConn(conn)
	return dbBase.GetDriver(ctx, conn, driverId)
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

func (db *DB) ListCloudDriverSync(ctx context.Context) (drivers []dao.Driver, err error) {
	conn := db.getConn()
	defer db.putConn(conn)
	return dbBase.ListCloudDriverSync(ctx, conn)
}

func (db *DB) ListLocalFileDriver(ctx context.Context, deviceId string) (drivers []dao.Driver, err error) {
	conn := db.getConn()
	defer db.putConn(conn)
	return dbBase.ListLocalFileDriver(ctx, conn, deviceId)
}

func (db *DB) GetDriverLocalFile(ctx context.Context, driverId uint64) (driver *dao.Driver, err error) {
	conn := db.getConn()
	defer db.putConn(conn)
	return dbBase.GetDriverLocalFile(ctx, conn, driverId)
}

func (db *DB) GetDriverDirCalculatedInfo(ctx context.Context, driverId uint64, filePath []string) (info dao.DirCalculatedInfo, err error) {
	conn := db.getConn()
	defer db.putConn(conn)
	return dbBase.GetDriverDirCalculatedInfo(ctx, conn, driverId, filePath)
}
