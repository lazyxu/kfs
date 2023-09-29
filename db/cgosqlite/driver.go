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
