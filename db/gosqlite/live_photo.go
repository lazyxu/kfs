package gosqlite

import (
	"context"

	"github.com/lazyxu/kfs/db/dbBase"
)

func (db *DB) SetLivpForMovAndHeicOrJpgAll(ctx context.Context) error {
	conn := db.getConn()
	defer db.putConn(conn)
	return dbBase.SetLivpForMovAndHeicOrJpgAll(ctx, conn)
}

func (db *DB) UpsertLivePhoto(ctx context.Context, movHash string, heicHash string, jpgHash string, livpHash string) error {
	conn := db.getConn()
	defer db.putConn(conn)
	return dbBase.UpsertLivePhoto(ctx, conn, movHash, heicHash, jpgHash, livpHash)
}

func (db *DB) ListLivePhotoNew(ctx context.Context) (hashList []string, err error) {
	conn := db.getConn()
	defer db.putConn(conn)
	return dbBase.ListLivePhotoNew(ctx, conn)
}

func (db *DB) SetLivpForMovAndHeicOrJpgInDirPath(ctx context.Context, driverId uint64, filePath []string) (err error) {
	conn := db.getConn()
	defer db.putConn(conn)
	return dbBase.SetLivpForMovAndHeicOrJpgInDirPath(ctx, conn, driverId, filePath)
}

func (db *DB) SetLivpForMovAndHeicOrJpgInDriver(ctx context.Context, driverId uint64) (err error) {
	conn := db.getConn()
	defer db.putConn(conn)
	return dbBase.SetLivpForMovAndHeicOrJpgInDriver(ctx, conn, driverId)
}

func (db *DB) ListLivePhotoAll(ctx context.Context) (hashList []string, err error) {
	conn := db.getConn()
	defer db.putConn(conn)
	return dbBase.ListLivePhotoAll(ctx, conn)
}

func (db *DB) GetLivePhotoByLivp(ctx context.Context, livpHash string) (movHash string, heicHash string, err error) {
	conn := db.getConn()
	defer db.putConn(conn)
	return dbBase.GetLivePhotoByLivp(ctx, conn, livpHash)
}
