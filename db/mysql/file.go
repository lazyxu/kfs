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

func (db *DB) GetDriverFile(ctx context.Context, driverName string, splitPath []string) (file dao.DriverFile, err error) {
	return dbBase.GetDriverFile(ctx, db.db, driverName, splitPath)
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
	return dbBase.UpsertDriverFileMysql(ctx, db.db, f)
}

func (db *DB) ListDriverFile(ctx context.Context, driverName string, filePath []string) (files []dao.DriverFile, err error) {
	return dbBase.ListDriverFile(ctx, db.db, driverName, filePath)
}

func (db *DB) InsertFile(ctx context.Context, hash string, size uint64) error {
	return dbBase.InsertFile(ctx, db.db, db, hash, size)
}

func (db *DB) InsertFileMd5(ctx context.Context, hash string, hashMd5 string) error {
	return dbBase.InsertFileMd5(ctx, db.db, db, hash, hashMd5)
}

func (db *DB) ListFileMd5(ctx context.Context, md5List []string) (m map[string]string, err error) {
	return dbBase.ListFileMd5(ctx, db.db, md5List)
}

func (db *DB) SumFileSize(ctx context.Context) (size uint64, err error) {
	return dbBase.SumFileSize(ctx, db.db)
}

func (db *DB) ListFile(ctx context.Context) (hashList []string, err error) {
	return dbBase.ListFile(ctx, db.db)
}

func (db *DB) ListDriverFileByHash(ctx context.Context, hash string) (files []dao.DriverFile, err error) {
	return dbBase.ListDriverFileByHash(ctx, db.db, hash)
}
