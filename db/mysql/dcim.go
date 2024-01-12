package mysql

import (
	"context"
	"github.com/lazyxu/kfs/dao"
	"github.com/lazyxu/kfs/db/dbBase"
)

func (db *DB) ListDCIMDriver(ctx context.Context) (drivers []dao.DCIMDriver, err error) {
	return dbBase.ListDCIMDriver(ctx, db.db)
}

func (db *DB) ListDCIMMediaType(ctx context.Context) (m map[string][]dao.Metadata, err error) {
	return dbBase.ListDCIMMediaType(ctx, db.db)
}

func (db *DB) ListDCIMLocation(ctx context.Context) (list []dao.Metadata, err error) {
	return dbBase.ListDCIMLocation(ctx, db.db)
}

func (db *DB) ListDCIMSearchType(ctx context.Context) (list []dao.DCIMSearchType, err error) {
	return dbBase.ListDCIMSearchType(ctx, db.db)
}

func (db *DB) ListDCIMSearchSuffix(ctx context.Context) (list []dao.DCIMSearchSuffix, err error) {
	return dbBase.ListDCIMSearchSuffix(ctx, db.db)
}
