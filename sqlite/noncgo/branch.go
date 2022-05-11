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

func (db *SqliteNonCgoDB) WriteBranch(ctx context.Context, branch Branch) error {
	_, err := db._db.ExecContext(ctx, `
	REPLACE INTO branch VALUES (?, ?, ?, ?, ?);
	`, branch.name, branch.description, branch.commitId, branch.size, branch.count)
	return err
}
