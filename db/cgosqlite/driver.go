package cgosqlite

import (
	"context"
	"github.com/lazyxu/kfs/dao"
	"github.com/lazyxu/kfs/db/dbBase"
)

func (db *DB) InsertDriver(ctx context.Context, driverName string, description string, typ string, accessToken string, refreshToken string) (exist bool, err error) {
	conn := db.getConn()
	defer db.putConn(conn)
	return dbBase.InsertDriver(ctx, conn, db, driverName, description, typ, accessToken, refreshToken)
}

func (db *DB) UpdateDriverSync(ctx context.Context, driverName string, sync bool, h int, m int, s int) error {
	conn := db.getConn()
	defer db.putConn(conn)
	return dbBase.UpdateDriverSync(ctx, conn, driverName, sync, h, m, s)
}

func (db *DB) DeleteDriver(ctx context.Context, driverName string) error {
	conn := db.getConn()
	defer db.putConn(conn)
	return dbBase.DeleteDriver(ctx, conn, driverName)
}

func (db *DB) ListDriver(ctx context.Context) (drivers []dao.Driver, err error) {
	conn := db.getConn()
	defer db.putConn(conn)
	return dbBase.ListDriver(ctx, conn)
}

func (db *DB) GetDriver(ctx context.Context, driverName string) (driver dao.Driver, err error) {
	conn := db.getConn()
	defer db.putConn(conn)
	return dbBase.GetDriver(ctx, conn, driverName)
}

func (db *DB) GetDriverFileSize(ctx context.Context, driverName string) (n uint64, err error) {
	conn := db.getConn()
	defer db.putConn(conn)
	return dbBase.GetDriverFileSize(ctx, conn, driverName)
}

func (db *DB) GetDriverFileCount(ctx context.Context, driverName string) (n uint64, err error) {
	conn := db.getConn()
	defer db.putConn(conn)
	return dbBase.GetDriverFileCount(ctx, conn, driverName)
}

func (db *DB) GetDriverDirCount(ctx context.Context, driverName string) (n uint64, err error) {
	conn := db.getConn()
	defer db.putConn(conn)
	return dbBase.GetDriverDirCount(ctx, conn, driverName)
}
