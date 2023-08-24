package cgosqlite

import (
	"context"
	"github.com/lazyxu/kfs/db/dbBase"
)

func (db *DB) UpsertLivePhoto(ctx context.Context, movHash string, heicHash string, jpgHash string) error {
	conn := db.getConn()
	defer db.putConn(conn)
	return dbBase.UpsertLivePhoto(ctx, conn, movHash, heicHash, jpgHash)
}
