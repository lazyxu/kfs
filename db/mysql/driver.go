package mysql

import (
	"context"
	"github.com/lazyxu/kfs/dao"
	"github.com/lazyxu/kfs/db/dbBase"
)

func (db *DB) InsertDriver(ctx context.Context, driverName string, description string, typ string, accessToken string, refreshToken string) (exist bool, err error) {
	return dbBase.InsertDriver(ctx, db.db, db, driverName, description, typ, accessToken, refreshToken)
}

func (db *DB) DeleteDriver(ctx context.Context, driverName string) error {
	return dbBase.DeleteDriver(ctx, db.db, driverName)
}

func (db *DB) ListDriver(ctx context.Context) (drivers []dao.Driver, err error) {
	return dbBase.ListDriver(ctx, db.db)
}
