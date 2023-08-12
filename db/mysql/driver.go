package mysql

import (
	"context"
	"github.com/lazyxu/kfs/dao"
	"github.com/lazyxu/kfs/db/dbBase"
)

func (db *DB) NewDriver(ctx context.Context, driverName string, description string) (exist bool, err error) {
	return dbBase.NewDriver(ctx, db.db, db, driverName, description)
}

func (db *DB) DeleteDriver(ctx context.Context, driverName string) error {
	return dbBase.DeleteDriver(ctx, db.db, driverName)
}

func (db *DB) DriverList(ctx context.Context) (drivers []dao.IDriver, err error) {
	return dbBase.DriverList(ctx, db.db)
}
