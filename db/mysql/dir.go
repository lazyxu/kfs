package mysql

import (
	"context"
	"github.com/lazyxu/kfs/dao"
	"github.com/lazyxu/kfs/db/dbBase"
)

func (db *DB) MaxBatchSize() int {
	return 65536
}

func (db *DB) WriteDir(ctx context.Context, dirItems []dao.DirItem) (dir dao.Dir, err error) {
	return dbBase.WriteDir(ctx, db.db, db, dirItems)
}

func (db *DB) RemoveDirItem(ctx context.Context, branchName string, splitPath []string) (commit dao.Commit, branch dao.Branch, err error) {
	return dbBase.RemoveDirItem(ctx, db.db, db, branchName, splitPath)
}
