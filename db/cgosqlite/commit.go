package cgosqlite

import (
	"context"
	"github.com/lazyxu/kfs/db/dbBase"

	"github.com/lazyxu/kfs/dao"
)

func (db *DB) WriteCommit(ctx context.Context, commit *dao.Commit) error {
	conn := db.getConn()
	defer db.putConn(conn)
	return db.InsertCommitWithTxOrDb(ctx, conn, commit)
}

func (db *DB) InsertCommitWithTxOrDb(ctx context.Context, txOrDb dbBase.TxOrDb, commit *dao.Commit) error {
	return dbBase.InsertCommitWithTxOrDbCgoSqlite(ctx, txOrDb, commit)
}
