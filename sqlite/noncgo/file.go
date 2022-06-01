package noncgo

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
)

type File struct {
	fileOrDir
}

func NewFile(hash string, size uint64) File {
	return File{fileOrDir{hash, size}}
}

func NewFileByBytes(bytes []byte) File {
	hash := sha256.New()
	hash.Write(bytes)
	return NewFile(hex.EncodeToString(hash.Sum(nil)), uint64(len(bytes)))
}

func (db *DB) WriteFile(ctx context.Context, file File) error {
	return db.writeFile(ctx, db._db, file)
}

func (db *DB) writeFile(ctx context.Context, txOrDb TxOrDb, file File) error {
	_, err := txOrDb.ExecContext(ctx, `
	INSERT INTO file VALUES (?, ?);
	`, file.hash, file.size)
	if err != nil {
		if isUniqueConstraintError(err) {
			return nil
		}
		return err
	}
	return err
}

func (db *DB) UpsertDirItem(ctx context.Context, branchName string, splitPath []string, item DirItem) (commit Commit, branch Branch, err error) {
	tx, err := db._db.Begin()
	if err != nil {
		return
	}
	defer func() {
		err = commitAndRollback(tx, err)
	}()
	if len(splitPath) == 0 {
		var dir Dir
		dir, err = NewDirFromDirItem(item)
		if err != nil {
			return
		}
		commit = NewCommit(dir, branchName, "")
		err = db.writeCommit(ctx, tx, &commit)
		if err != nil {
			return
		}
		branch = NewBranch(branchName, commit, dir)
		err = db.writeBranch(ctx, tx, branch)
		return
	}
	return db.updateDirItem(ctx, tx, branchName, splitPath, func(dirItemsList [][]DirItem) error {
		i := len(dirItemsList) - 1
		item.Name = splitPath[i]
		find := false
		for j, dirItem := range dirItemsList[i] {
			if dirItem.Name == splitPath[i] {
				dirItemsList[i][j] = item // update
				find = true
				break
			}
		}
		if !find {
			dirItemsList[i] = append(dirItemsList[i], item) // insert
		}
		return nil
	})
}
