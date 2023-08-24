package mysql

import (
	"context"
	"github.com/lazyxu/kfs/db/dbBase"
)

func (db *DB) UpsertLivePhoto(ctx context.Context, movHash string, heicHash string, jpgHash string) error {
	return dbBase.UpsertLivePhoto(ctx, db.db, movHash, heicHash, jpgHash)
}
