package cgosqlite

func (db *DB) IsUniqueConstraintError(err error) bool {
	if err == nil {
		return false
	}
	return false
}
