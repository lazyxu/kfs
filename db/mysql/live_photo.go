package mysql

import (
	"context"

	"github.com/lazyxu/kfs/db/dbBase"
)

func (db *DB) SetLivpForMovAndHeicOrJpgAll(ctx context.Context) error {
	return dbBase.SetLivpForMovAndHeicOrJpgAll(ctx, db.db)
}

func (db *DB) UpsertLivePhoto(ctx context.Context, movHash string, heicHash string, jpgHash string, livpHash string) error {
	return dbBase.UpsertLivePhoto(ctx, db.db, movHash, heicHash, jpgHash, livpHash)
}

func (db *DB) ListLivePhotoNew(ctx context.Context) (hashList []string, err error) {
	return dbBase.ListLivePhotoNew(ctx, db.db)
}

func (db *DB) SetLivpForMovAndHeicOrJpgInDirPath(ctx context.Context, driverId uint64, filePath []string) (err error) {
	return dbBase.SetLivpForMovAndHeicOrJpgInDirPath(ctx, db.db, driverId, filePath)
}

func (db *DB) SetLivpForMovAndHeicOrJpgInDriver(ctx context.Context, driverId uint64) (err error) {
	return dbBase.SetLivpForMovAndHeicOrJpgInDriver(ctx, db.db, driverId)
}

func (db *DB) ListLivePhotoAll(ctx context.Context) (hashList []string, err error) {
	return dbBase.ListLivePhotoAll(ctx, db.db)
}

func (db *DB) GetLivePhotoByLivp(ctx context.Context, livpHash string) (movHash string, heicHash string, err error) {
	return dbBase.GetLivePhotoByLivp(ctx, db.db, livpHash)
}
