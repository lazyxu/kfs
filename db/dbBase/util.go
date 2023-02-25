package dbBase

import "database/sql"

func CommitAndRollback(tx *sql.Tx, err error) error {
	if err != nil {
		err1 := tx.Rollback()
		if err1 != nil {
			return err1
		}
		return err
	}
	err = tx.Commit()
	if err == nil {
		return nil
	}
	err1 := tx.Rollback()
	if err1 != nil {
		return err1
	}
	return err
}
