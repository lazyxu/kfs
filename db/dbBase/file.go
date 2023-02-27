package dbBase

import (
	"context"
	"database/sql"
	"errors"
	"github.com/lazyxu/kfs/dao"
	"os"
	"strings"
)

func WriteFileWithTxOrDb(ctx context.Context, txOrDb TxOrDb, db DbImpl, file dao.File) error {
	_, err := txOrDb.ExecContext(ctx, `
	INSERT INTO _file VALUES (?, ?);
	`, file.Hash(), file.Size())
	if err != nil {
		if db.IsUniqueConstraintError(err) {
			return nil
		}
		return err
	}
	return err
}

func GetFile(ctx context.Context, conn *sql.DB, branchName string, splitPath []string) (dirItem dao.DirItem, err error) {
	tx, err := conn.Begin()
	if err != nil {
		return
	}
	defer func() {
		err = CommitAndRollback(tx, err)
	}()
	if len(splitPath) == 0 {
		err = errors.New("/: Is a directory")
		return
	}
	hash, err := getBranchCommitHash(ctx, tx, branchName)
	if err != nil {
		return
	}
	for i := range splitPath[:len(splitPath)-1] {
		hash, err = getDirItemHash(ctx, tx, hash, splitPath, i)
		if err != nil {
			return
		}
	}
	dirItem, err = getDirItem(ctx, tx, hash, splitPath, len(splitPath)-1)
	if err != nil {
		return
	}
	if os.FileMode(dirItem.Mode).IsDir() {
		err = errors.New("/" + strings.Join(splitPath, "/") + ": Is a directory")
		return
	}
	return
}

func UpsertDirItem(ctx context.Context, conn *sql.DB, db DbImpl, branchName string, splitPath []string, item dao.DirItem) (commit dao.Commit, branch dao.Branch, err error) {
	tx, err := conn.Begin()
	if err != nil {
		return
	}
	defer func() {
		err = CommitAndRollback(tx, err)
	}()
	if len(splitPath) == 0 {
		var dir dao.Dir
		dir, err = dao.NewDirFromDirItem(item)
		if err != nil {
			return
		}
		commit = dao.NewCommit(dir, branchName, "")
		err = db.InsertCommitWithTxOrDb(ctx, tx, &commit)
		if err != nil {
			return
		}
		branch = dao.NewBranch(branchName, commit, dir)
		err = db.UpsertBranchWithTxOrDb(ctx, tx, branch)
		if err != nil {
			return
		}
		return
	}
	return updateDirItemWithTx(ctx, tx, db, branchName, splitPath, func(dirItemsList [][]dao.DirItem) ([]dao.DirItem, error) {
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

func UpsertDirItems(ctx context.Context, conn *sql.DB, db DbImpl, branchName string, splitPath []string, items []dao.DirItem) (commit dao.Commit, branch dao.Branch, err error) {
	tx, err := conn.Begin()
	if err != nil {
		return
	}
	defer func() {
		err = CommitAndRollback(tx, err)
	}()
	return updateDirItemsWithTx(ctx, tx, db, branchName, splitPath, func(dirItems *[]dao.DirItem) ([]dao.DirItem, error) {
		for _, item := range items {
			find := false
			for j, dirItem := range *dirItems {
				if dirItem.Name == item.Name {
					(*dirItems)[j] = item // update
					find = true
					break
				}
			}
			if !find {
				*dirItems = append(*dirItems, item) // insert
			}
		}
		return items, nil
	})
}

func GetFileHashMode(ctx context.Context, conn *sql.DB, branchName string, splitPath []string) (hash string, mode os.FileMode, err error) {
	tx, err := conn.Begin()
	if err != nil {
		return
	}
	defer func() {
		err = CommitAndRollback(tx, err)
	}()
	hash, err = getBranchCommitHash(ctx, tx, branchName)
	if err != nil {
		return
	}
	if len(splitPath) == 0 {
		return hash, os.ModeDir | os.ModePerm, nil
	}
	for i := range splitPath[:len(splitPath)-1] {
		hash, err = getDirItemHash(ctx, tx, hash, splitPath, i)
		if err != nil {
			return
		}
	}
	hash, m, err := getDirItemHashMode(ctx, tx, hash, splitPath, len(splitPath)-1)
	mode = os.FileMode(m)
	return
}
