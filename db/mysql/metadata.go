package mysql

import (
	"context"
	"github.com/lazyxu/kfs/dao"
	"github.com/lazyxu/kfs/db/dbBase"
)

func (db *DB) InsertDCIMMetadataTime(ctx context.Context, hash string, t int64) (exist bool, err error) {
	return dbBase.InsertDCIMMetadataTime(ctx, db.db, db, hash, t)
}

func (db *DB) GetEarliestCrated(ctx context.Context, hash string) int64 {
	return dbBase.GetEarliestCrated(ctx, db.db, db, hash)
}

func (db *DB) ListMetadataTime(ctx context.Context) (list []dao.Metadata, err error) {
	return dbBase.ListMetadataTime(ctx, db.db)
}
