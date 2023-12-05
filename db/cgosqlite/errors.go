package cgosqlite

import "github.com/mattn/go-sqlite3"

func (db *DB) IsUniqueConstraintError(err error) bool {
	if err == nil {
		return false
	}
	if e, ok := err.(sqlite3.Error); ok {
		return e.Code == sqlite3.ErrConstraint
	}
	return false
}
