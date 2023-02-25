package mysql

import (
	"context"
	"github.com/lazyxu/kfs/db/dbBase"

	"github.com/lazyxu/kfs/dao"
)

func (db *DB) WriteCommit(ctx context.Context, commit *dao.Commit) error {
	return db.InsertCommitWithTxOrDb(ctx, db.db, commit)
}

func (db *DB) InsertCommitWithTxOrDb(ctx context.Context, txOrDb dbBase.TxOrDb, commit *dao.Commit) error {
	return dbBase.InsertCommitWithTxOrDb(ctx, txOrDb, commit)
}
