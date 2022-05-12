package noncgo

import "context"

type Branch struct {
	name        string
	description string
	commitId    int64
	size        uint64
	count       uint64
}

func NewBranch(name string, description string, commit Commit, dir Dir) Branch {
	return Branch{name, description, commit.id, dir.size, dir.count}
}

func (db *DB) WriteBranch(ctx context.Context, branch Branch) error {
	return db.writeBranch(ctx, db._db, branch)
}

func (db *DB) writeBranch(ctx context.Context, txOrDb TxOrDb, branch Branch) error {
	_, err := txOrDb.ExecContext(ctx, `
	REPLACE INTO branch VALUES (?, ?, ?, ?, ?);
	`, branch.name, branch.description, branch.commitId, branch.size, branch.count)
	return err
}
