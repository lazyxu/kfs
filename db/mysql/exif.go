package mysql

import (
	"context"
	"github.com/lazyxu/kfs/dao"
	"github.com/lazyxu/kfs/db/dbBase"
)

func (db *DB) InsertHeightWidth(ctx context.Context, hash string, hw dao.HeightWidth) error {
	return dbBase.InsertHeightWidth(ctx, db.db, db, hash, hw)
}

func (db *DB) InsertNullVideoMetadata(ctx context.Context, hash string) (exist bool, err error) {
	return dbBase.InsertNullVideoMetadata(ctx, db.db, db, hash)
}

func (db *DB) InsertVideoMetadata(ctx context.Context, hash string, m dao.VideoMetadata) (exist bool, err error) {
	return dbBase.InsertVideoMetadata(ctx, db.db, db, hash, m)
}

func (db *DB) InsertNullExif(ctx context.Context, hash string) (exist bool, err error) {
	return dbBase.InsertNullExif(ctx, db.db, db, hash)
}

func (db *DB) InsertExif(ctx context.Context, hash string, e dao.Exif) (exist bool, err error) {
	return dbBase.InsertExif(ctx, db.db, db, hash, e)
}

func (db *DB) ListExpectExif(ctx context.Context) (hashList []string, err error) {
	return dbBase.ListExpectExif(ctx, db.db)
}

func (db *DB) ListExpectExifCb(ctx context.Context, cb func(hash string)) (err error) {
	return dbBase.ListExpectExifCb(ctx, db.db, cb)
}

func (db *DB) ListExif(ctx context.Context) (exifMap map[string]dao.Exif, err error) {
	return dbBase.ListExif(ctx, db.db)
}

func (db *DB) ListMetadata(ctx context.Context) (list []dao.Metadata, err error) {
	return dbBase.ListMetadata(ctx, db.db)
}

func (db *DB) GetMetadata(ctx context.Context, hash string) (metadata dao.Metadata, err error) {
	return dbBase.GetMetadata(ctx, db.db, hash)
}
