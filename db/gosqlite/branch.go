package gosqlite

import (
	"context"
	"github.com/lazyxu/kfs/db/dbBase"

	"github.com/lazyxu/kfs/dao"
)

func (db *DB) ResetBranch(ctx context.Context, branchName string) (err error) {
	conn := db.getConn()
	defer db.putConn(conn)
	return dbBase.ResetBranch(ctx, conn, db, branchName)
}

func (db *DB) WriteBranch(ctx context.Context, branch dao.Branch) error {
	conn := db.getConn()
	defer db.putConn(conn)
	return db.UpsertBranchWithTxOrDb(ctx, conn, branch)
}

func (db *DB) NewBranch(ctx context.Context, branchName string) (exist bool, err error) {
	conn := db.getConn()
	defer db.putConn(conn)
	return dbBase.NewBranch(ctx, conn, db, branchName)
}

func (db *DB) BranchInfo(ctx context.Context, branchName string) (branch dao.Branch, err error) {
	conn := db.getConn()
	defer db.putConn(conn)
	return dbBase.BranchInfo(ctx, conn, branchName)
}

func (db *DB) BranchList(ctx context.Context) (branches []dao.IBranch, err error) {
	conn := db.getConn()
	defer db.putConn(conn)
	return dbBase.BranchList(ctx, conn)
}

func (db *DB) UpsertBranchWithTxOrDb(ctx context.Context, txOrDb dbBase.TxOrDb, branch dao.Branch) error {
	return dbBase.UpsertBranchWithTxOrDb(ctx, txOrDb, branch)
}
