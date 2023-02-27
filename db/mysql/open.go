package mysql

import (
	"context"
	"github.com/lazyxu/kfs/db/dbBase"
	"os"

	"github.com/lazyxu/kfs/dao"
)

func (db *DB) Open(ctx context.Context, branchName string, splitPath []string) (hash string, mode os.FileMode, dirItems []dao.DirItem, err error) {
	return dbBase.Open(ctx, db.db, branchName, splitPath)
}

func (db *DB) Open2(ctx context.Context, branchName string, splitPath []string) (dirItem dao.DirItem, dirItems []dao.DirItem, err error) {
	return dbBase.Open2(ctx, db.db, branchName, splitPath)
}
