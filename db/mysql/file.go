package mysql

import (
	"context"
	"github.com/lazyxu/kfs/db/dbBase"
	"os"

	"github.com/lazyxu/kfs/dao"
)

func (db *DB) WriteFile(ctx context.Context, file dao.File) error {
	return dbBase.WriteFileWithTxOrDb(ctx, db.db, db, file)
}

func (db *DB) GetFile(ctx context.Context, driverName string, splitPath []string) (file dao.DriverFile, err error) {
	return dbBase.GetFile(ctx, db.db, driverName, splitPath)
}

func (db *DB) UpsertDirItem(ctx context.Context, branchName string, splitPath []string, item dao.DirItem) (commit dao.Commit, branch dao.Branch, err error) {
	return dbBase.UpsertDirItem(ctx, db.db, db, branchName, splitPath, item)
}

func (db *DB) UpsertDirItems(ctx context.Context, branchName string, splitPath []string, items []dao.DirItem) (commit dao.Commit, branch dao.Branch, err error) {
	return dbBase.UpsertDirItems(ctx, db.db, db, branchName, splitPath, items)
}

func (db *DB) GetFileHashMode(ctx context.Context, branchName string, splitPath []string) (hash string, mode os.FileMode, err error) {
	return dbBase.GetFileHashMode(ctx, db.db, branchName, splitPath)
}

func (db *DB) UpsertDriverFile(ctx context.Context, f dao.DriverFile) error {
	return dbBase.UpsertDriverFile(ctx, db.db, f)
}

func (db *DB) ListV2(ctx context.Context, driverName string, filePath []string) (files []dao.DriverFile, err error) {
	return dbBase.ListV2(ctx, db.db, driverName, filePath)
}
