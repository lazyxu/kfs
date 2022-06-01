package noncgo

import (
	"context"
	"database/sql"

	"modernc.org/sqlite"
)

func commitAndRollback(tx *sql.Tx, err error) error {
	if err != nil {
		err1 := tx.Rollback()
		if err1 != nil {
			panic(err1) // should not happen
		}
		if e, ok := err.(*sqlite.Error); ok {
			if e.Code() == 5 {
				return nil
			}
			// constraint failed: UNIQUE constraint failed: hash.hashval (1555)
			if e.Code() == 1555 {
				return nil
			}
		}
		return err
	}
	err = tx.Commit()
	if err == nil {
		return nil
	}
	err1 := tx.Rollback()
	if err1 != nil {
		panic(err1) // should not happen
	}
	return err
}

type FileOrDir interface {
	Hash() string
	Size() uint64
	Count() uint64
	TotalCount() uint64
}

type fileOrDir struct {
	hash string
	size uint64
}

func (i fileOrDir) Hash() string {
	return i.hash
}

func (i fileOrDir) Size() uint64 {
	return i.size
}

func (i fileOrDir) Count() uint64 {
	return 1
}

func (i fileOrDir) TotalCount() uint64 {
	return 1
}

func (db *DB) count(ctx context.Context, tableName string) (int, error) {
	conn := db.getConn()
	defer db.putConn(conn)
	rows, err := conn.QueryContext(ctx, "SELECT COUNT(1) FROM "+tableName+";")
	if err != nil {
		return 0, err
	}
	if err = rows.Err(); err != nil {
		return 0, err
	}
	defer rows.Close()
	if rows.Next() {
		var i int
		if err = rows.Scan(&i); err != nil {
			return 0, err
		}
		return i, nil
	}
	panic("internal error when get " + tableName + " count")
}

func (db *DB) FileCount(ctx context.Context) (int, error) {
	return db.count(ctx, "file")
}

func (db *DB) DirCount(ctx context.Context) (int, error) {
	return db.count(ctx, "dir")
}

func (db *DB) DirItemCount(ctx context.Context) (int, error) {
	return db.count(ctx, "dirItem")
}

func (db *DB) BranchCount(ctx context.Context) (int, error) {
	return db.count(ctx, "branch")
}
