package noncgo

import "context"

type Branch struct {
	name        string
	description string
	hash        string
	size        uint64
	count       uint64
}

func NewBranch(name string, description string, dir Dir) Branch {
	return Branch{name, description, dir.hash, dir.size, dir.count}
}

func (db *SqliteNonCgoDB) WriteBranch(ctx context.Context, branch Branch) error {
	_, err := db._db.ExecContext(ctx, `
	INSERT INTO branch VALUES (?, ?, ?, ?, ?);
	`, branch.name, branch.description, branch.hash, branch.size, branch.count)
	return err
}
