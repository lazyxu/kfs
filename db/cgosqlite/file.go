package cgosqlite

import (
	"context"
	"github.com/lazyxu/kfs/db/dbBase"
	"os"

	"github.com/lazyxu/kfs/dao"
)

func (db *DB) WriteFile(ctx context.Context, file dao.File) error {
	conn := db.getConn()
	defer db.putConn(conn)
	return dbBase.WriteFileWithTxOrDb(ctx, conn, db, file)
}

func (db *DB) GetFile(ctx context.Context, branchName string, splitPath []string) (dirItem dao.DirItem, err error) {
	conn := db.getConn()
	defer db.putConn(conn)
	return dbBase.GetFile(ctx, conn, branchName, splitPath)
}

func (db *DB) UpsertDirItem(ctx context.Context, branchName string, splitPath []string, item dao.DirItem) (commit dao.Commit, branch dao.Branch, err error) {
	conn := db.getConn()
	defer db.putConn(conn)
	return dbBase.UpsertDirItem(ctx, conn, db, branchName, splitPath, item)
}

func (db *DB) UpsertDirItems(ctx context.Context, branchName string, splitPath []string, items []dao.DirItem) (commit dao.Commit, branch dao.Branch, err error) {
	conn := db.getConn()
	defer db.putConn(conn)
	return dbBase.UpsertDirItems(ctx, conn, db, branchName, splitPath, items)
}

func (db *DB) GetFileHashMode(ctx context.Context, branchName string, splitPath []string) (hash string, mode os.FileMode, err error) {
	conn := db.getConn()
	defer db.putConn(conn)
	return dbBase.GetFileHashMode(ctx, conn, branchName, splitPath)
}

func (db *DB) UpsertDriverFile(ctx context.Context, f dao.DriverFile) error {
	conn := db.getConn()
	defer db.putConn(conn)
	return dbBase.UpsertDriverFile(ctx, conn, f)
}

func (db *DB) ListV2(ctx context.Context, driverName string, filePath []string) (files []dao.DriverFile, err error) {
	conn := db.getConn()
	defer db.putConn(conn)
	return dbBase.ListV2(ctx, conn, driverName, filePath)
}
