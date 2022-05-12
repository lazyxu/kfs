package noncgo

import "context"

type FileOrDir interface {
	Hash() string
	Size() uint64
	Count() uint64
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

func (db *DB) count(ctx context.Context, tableName string) (int, error) {
	rows, err := db._db.QueryContext(ctx, "SELECT COUNT(1) FROM "+tableName+";")
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
