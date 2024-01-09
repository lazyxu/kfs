package gosqlite

import (
	"context"
	"github.com/lazyxu/kfs/dao"
	"github.com/lazyxu/kfs/db/dbBase"
)

func (db *DB) InsertHeightWidth(ctx context.Context, hash string, hw dao.HeightWidth) error {
	conn := db.getConn()
	defer db.putConn(conn)
	return dbBase.InsertHeightWidth(ctx, conn, db, hash, hw)
}
func (db *DB) InsertNullVideoMetadata(ctx context.Context, hash string) (exist bool, err error) {
	conn := db.getConn()
	defer db.putConn(conn)
	return dbBase.InsertNullVideoMetadata(ctx, conn, db, hash)
}

func (db *DB) InsertVideoMetadata(ctx context.Context, hash string, m dao.VideoMetadata) (exist bool, err error) {
	conn := db.getConn()
	defer db.putConn(conn)
	return dbBase.InsertVideoMetadata(ctx, conn, db, hash, m)
}

func (db *DB) InsertNullExif(ctx context.Context, hash string) (exist bool, err error) {
	conn := db.getConn()
	defer db.putConn(conn)
	return dbBase.InsertNullExif(ctx, conn, db, hash)
}

func (db *DB) InsertExif(ctx context.Context, hash string, e dao.Exif) (exist bool, err error) {
	conn := db.getConn()
	defer db.putConn(conn)
	return dbBase.InsertExif(ctx, conn, db, hash, e)
}

func (db *DB) ListExpectExif(ctx context.Context) (hashList []string, err error) {
	conn := db.getConn()
	defer db.putConn(conn)
	return dbBase.ListExpectExif(ctx, conn)
}

func (db *DB) ListExpectExifCb(ctx context.Context, cb func(hash string)) (err error) {
	conn := db.getConn()
	defer db.putConn(conn)
	return dbBase.ListExpectExifCb(ctx, conn, cb)
}

func (db *DB) ListExif(ctx context.Context) (exifMap map[string]dao.Exif, err error) {
	conn := db.getConn()
	defer db.putConn(conn)
	return dbBase.ListExif(ctx, conn)
}

func (db *DB) ListMetadata(ctx context.Context) (list []dao.Metadata, err error) {
	conn := db.getConn()
	defer db.putConn(conn)
	return dbBase.ListMetadata(ctx, conn)
}

func (db *DB) GetMetadata(ctx context.Context, hash string) (metadata dao.Metadata, err error) {
	conn := db.getConn()
	defer db.putConn(conn)
	return dbBase.GetMetadata(ctx, conn, hash)
}
