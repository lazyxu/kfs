package gosqlite

import (
	"context"
	"github.com/lazyxu/kfs/dao"
	"github.com/lazyxu/kfs/db/dbBase"
)

func (db *DB) NewDriver(ctx context.Context, driverName string, description string) (exist bool, err error) {
	conn := db.getConn()
	defer db.putConn(conn)
	return dbBase.NewDriver(ctx, conn, db, driverName, description)
}

func (db *DB) DeleteDriver(ctx context.Context, driverName string) error {
	conn := db.getConn()
	defer db.putConn(conn)
	return dbBase.DeleteDriver(ctx, conn, driverName)
}

func (db *DB) DriverList(ctx context.Context) (drivers []dao.IDriver, err error) {
	conn := db.getConn()
	defer db.putConn(conn)
	return dbBase.DriverList(ctx, conn)
}
