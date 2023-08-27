package mysql

import (
	"context"
	"github.com/lazyxu/kfs/dao"
	"github.com/lazyxu/kfs/db/dbBase"
)

func (db *DB) InsertFileType(ctx context.Context, hash string, t dao.FileType) (exist bool, err error) {
	return dbBase.InsertFileType(ctx, db.db, db, hash, t)
}

func (db *DB) ListExpectFileType(ctx context.Context) (hashList []string, err error) {
	return dbBase.ListExpectFileType(ctx, db.db)
}

func (db *DB) GetFileType(ctx context.Context, hash string) (fileType *dao.FileType, err error) {
	return dbBase.GetFileType(ctx, db.db, hash)
}
