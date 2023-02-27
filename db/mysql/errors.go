package mysql

import (
	"github.com/go-sql-driver/mysql"
)

func (db *DB) IsUniqueConstraintError(err error) bool {
	if err == nil {
		return false
	}
	if e, ok := err.(*mysql.MySQLError); ok {
		return e.Number == 1062
	}
	return false
}
