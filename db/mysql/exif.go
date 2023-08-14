package mysql

import (
	"context"
	"github.com/lazyxu/kfs/dao"
	"github.com/lazyxu/kfs/db/dbBase"
)

func (db *DB) InsertNullExif(ctx context.Context, hash string) (exist bool, err error) {
	return dbBase.InsertNullExif(ctx, db.db, db, hash)
}

func (db *DB) InsertExif(ctx context.Context, hash string, e dao.ExifData) (exist bool, err error) {
	return dbBase.InsertExif(ctx, db.db, db, hash, e)
}

func (db *DB) ListExpectExif(ctx context.Context) (hashList []string, err error) {
	return dbBase.ListExpectExif(ctx, db.db)
}

func (db *DB) ListExpectExifCb(ctx context.Context, cb func(hash string)) (err error) {
	return dbBase.ListExpectExifCb(ctx, db.db, cb)
}
