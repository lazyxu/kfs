package mysql

import (
	"context"
	"github.com/lazyxu/kfs/db/dbBase"

	"github.com/lazyxu/kfs/dao"
)

func (db *DB) List(ctx context.Context, branchName string, splitPath []string) (dirItems []dao.DirItem, err error) {
	return dbBase.List(ctx, db.db, branchName, splitPath)
}

func (db *DB) ListByHash(ctx context.Context, hash string) (dirItems []dao.DirItem, err error) {
	return dbBase.ListByHash(ctx, db.db, hash)
}
