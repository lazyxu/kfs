package gosqlite

import (
	"context"
	"os"

	"github.com/lazyxu/kfs/dao"
)

func (db *DB) WriteFile(ctx context.Context, file dao.File) error {
	conn := db.getConn()
	defer db.putConn(conn)
	return db.writeFile(ctx, conn, file)
}

func (db *DB) writeFile(ctx context.Context, txOrDb TxOrDb, file dao.File) error {
	_, err := txOrDb.ExecContext(ctx, `
	INSERT INTO file VALUES (?, ?);
	`, file.Hash(), file.Size())
	if err != nil {
		if isUniqueConstraintError(err) {
			return nil
		}
		return err
	}
	return err
}

func (db *DB) UpsertDirItem(ctx context.Context, branchName string, splitPath []string, item dao.DirItem) (commit dao.Commit, branch dao.Branch, err error) {
	conn := db.getConn()
	defer db.putConn(conn)
	tx, err := conn.Begin()
	if err != nil {
		return
	}
	defer func() {
		err = commitAndRollback(tx, err)
	}()
	if len(splitPath) == 0 {
		var dir dao.Dir
		dir, err = dao.NewDirFromDirItem(item)
		if err != nil {
			return
		}
		commit = dao.NewCommit(dir, branchName, "")
		err = db.writeCommit(ctx, tx, &commit)
		if err != nil {
			return
		}
		branch = dao.NewBranch(branchName, commit, dir)
		err = db.writeBranch(ctx, tx, branch)
		return
	}
	return db.updateDirItem(ctx, tx, branchName, splitPath, func(dirItemsList [][]dao.DirItem) ([]dao.DirItem, error) {
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
		return []dao.DirItem{item}, nil
	})
}

func (db *DB) GetFileHashMode(ctx context.Context, branchName string, splitPath []string) (hash string, mode os.FileMode, err error) {
	conn := db.getConn()
	defer db.putConn(conn)
	tx, err := conn.Begin()
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
	hash, err = db.getBranchCommitHash(ctx, tx, branchName)
	if err != nil {
		return
	}
	if len(splitPath) == 0 {
		return hash, os.ModeDir | os.ModePerm, nil
	}
	for i := range splitPath[:len(splitPath)-1] {
		hash, err = db.getDirItemHash(ctx, tx, hash, splitPath, i)
		if err != nil {
			return
		}
	}
	hash, m, err := db.getDirItemHashMode(ctx, tx, hash, splitPath, len(splitPath)-1)
	mode = os.FileMode(m)
	return
}
