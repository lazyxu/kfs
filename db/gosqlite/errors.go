package gosqlite

import (
	"modernc.org/sqlite"
)

func isUniqueConstraintError(err error) bool {
	if err == nil {
		return false
	}
	if e, ok := err.(*sqlite.Error); ok {
		return e.Code() == 1555
	}
	return false
}
