package gosqlite

import (
	"context"
	"github.com/lazyxu/kfs/dao"
	"github.com/lazyxu/kfs/db/dbBase"
	"os"
)

func (db *DB) WriteFile(ctx context.Context, file dao.File) error {
	conn := db.getConn()
	defer db.putConn(conn)
	return dbBase.WriteFileWithTxOrDb(ctx, conn, db, file)
}

func (db *DB) GetDriverFile(ctx context.Context, driverId uint64, splitPath []string) (file dao.DriverFile, err error) {
	conn := db.getConn()
	defer db.putConn(conn)
	return dbBase.GetDriverFile(ctx, conn, driverId, splitPath)
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

func (db *DB) ListDriverFile(ctx context.Context, driverId uint64, filePath []string) (files []dao.DriverFile, err error) {
	conn := db.getConn()
	defer db.putConn(conn)
	return dbBase.ListDriverFile(ctx, conn, driverId, filePath)
}

func (db *DB) InsertFile(ctx context.Context, hash string, size uint64) error {
	conn := db.getConn()
	defer db.putConn(conn)
	return dbBase.InsertFile(ctx, conn, db, hash, size)
}

func (db *DB) InsertFileMd5(ctx context.Context, hash string, hashMd5 string) error {
	conn := db.getConn()
	defer db.putConn(conn)
	return dbBase.InsertFileMd5(ctx, conn, db, hash, hashMd5)
}

func (db *DB) ListFileMd5(ctx context.Context, md5List []string) (m map[string]string, err error) {
	conn := db.getConn()
	defer db.putConn(conn)
	return dbBase.ListFileMd5(ctx, conn, md5List)
}

func (db *DB) SumFileSize(ctx context.Context) (size uint64, err error) {
	conn := db.getConn()
	defer db.putConn(conn)
	return dbBase.SumFileSize(ctx, conn)
}

func (db *DB) ListFile(ctx context.Context) (hashList []string, err error) {
	conn := db.getConn()
	defer db.putConn(conn)
	return dbBase.ListFile(ctx, conn)
}

func (db *DB) ListDriverFileByHash(ctx context.Context, hash string) (files []dao.DriverFile, err error) {
	conn := db.getConn()
	defer db.putConn(conn)
	return dbBase.ListDriverFileByHash(ctx, conn, hash)
}
