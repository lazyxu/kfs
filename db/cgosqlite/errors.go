package cgosqlite

func isUniqueConstraintError(err error) bool {
	if err == nil {
		return false
	}
	return false
}

func (db *DB) IsUniqueConstraintError(err error) bool {
	if err == nil {
		return false
	}
	return false
}
