package mysql

import (
	"context"
	"github.com/lazyxu/kfs/db/dbBase"

	"github.com/lazyxu/kfs/dao"
)

func (db *DB) ResetBranch(ctx context.Context, branchName string) (err error) {
	return dbBase.ResetBranch(ctx, db.db, db, branchName)
}

func (db *DB) WriteBranch(ctx context.Context, branch dao.Branch) error {
	return db.UpsertBranchWithTxOrDb(ctx, db.db, branch)
}

func (db *DB) NewBranch(ctx context.Context, branchName string) (exist bool, err error) {
	return dbBase.NewBranch(ctx, db.db, db, branchName)
}

func (db *DB) DeleteBranch(ctx context.Context, branchName string) error {
	return dbBase.DeleteBranch(ctx, db.db, branchName)
}

func (db *DB) BranchInfo(ctx context.Context, branchName string) (branch dao.Branch, err error) {
	return dbBase.BranchInfo(ctx, db.db, branchName)
}

func (db *DB) BranchList(ctx context.Context) (branches []dao.IBranch, err error) {
	return dbBase.BranchList(ctx, db.db)
}

func (db *DB) UpsertBranchWithTxOrDb(ctx context.Context, txOrDb dbBase.TxOrDb, branch dao.Branch) error {
	return dbBase.UpsertBranchWithTxOrDbMysql(ctx, txOrDb, branch)
}

func (db *DB) InsertBranchWithTxOrDb(ctx context.Context, txOrDb dbBase.TxOrDb, branch dao.Branch) error {
	return dbBase.InsertBranchWithTxOrDb(ctx, txOrDb, branch)
}
