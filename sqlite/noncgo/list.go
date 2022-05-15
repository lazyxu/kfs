package noncgo

import "context"

func (db *DB) List(ctx context.Context, branchName string, splitPath []string) (dirItems []DirItem, err error) {
	tx, err := db._db.Begin()
	if err != nil {
		return
	}
	defer func() {
		if err == nil {
			err = tx.Commit()
			if err != nil {
				err1 := tx.Rollback()
				if err1 != nil {
					panic(err1) // should not happen
				}
				return
			}
		}
	}()
	hash, err := db.getBranchCommitHash(ctx, tx, branchName)
	if err != nil {
		return
	}
	for i := range splitPath {
		hash, err = db.getDirItemHash(ctx, tx, hash, splitPath, i)
		if err != nil {
			return
		}
	}
	dirItems, err = db.getDirItems(ctx, tx, hash)
	if err != nil {
		return
	}
	return
}
