package dbBase

import (
	"context"
	"database/sql"
	"errors"
	"github.com/lazyxu/kfs/dao"
	"strings"
)

func getDirItemHash(ctx context.Context, tx *sql.Tx, hash string, splitPath []string, i int) (itemHash string, err error) {
	rows, err := tx.QueryContext(ctx, `
		SELECT itemHash FROM _dirItem WHERE Hash=? and itemName=?
	`, hash, splitPath[i])
	if err != nil {
		return
	}
	defer rows.Close()
	if !rows.Next() {
		return "", errors.New("no such file or dir: /" + strings.Join(splitPath, "/"))
	}
	err = rows.Scan(&itemHash)
	return
}

func getDirItemHashMode(ctx context.Context, tx *sql.Tx, hash string, splitPath []string, i int) (itemHash string, itemMode uint64, err error) {
	rows, err := tx.QueryContext(ctx, `
		SELECT itemHash, itemMode FROM _dirItem WHERE Hash=? and itemName=?
	`, hash, splitPath[i])
	if err != nil {
		return
	}
	defer rows.Close()
	if !rows.Next() {
		return "", 0, errors.New("no such file or dir: /" + strings.Join(splitPath, "/"))
	}
	err = rows.Scan(&itemHash, &itemMode)
	return
}

func getDirItem(ctx context.Context, tx *sql.Tx, hash string, splitPath []string, i int) (dirItem dao.DirItem, err error) {
	rows, err := tx.QueryContext(ctx, `
		SELECT itemHash,
			itemName,
			itemMode,
			itemSize,
			itemCount,
			itemTotalCount,
			itemCreateTime,
			itemModifyTime,
			itemChangeTime,
			itemAccessTime
		FROM _dirItem WHERE Hash=? and itemName=?
	`, hash, splitPath[i])
	if err != nil {
		return
	}
	defer rows.Close()
	if !rows.Next() {
		err = errors.New("no such file or dir: /" + strings.Join(splitPath, "/"))
		return
	}
	err = rows.Scan(
		&dirItem.Hash,
		&dirItem.Name,
		&dirItem.Mode,
		&dirItem.Size,
		&dirItem.Count,
		&dirItem.TotalCount,
		&dirItem.CreateTime,
		&dirItem.ModifyTime,
		&dirItem.ChangeTime,
		&dirItem.AccessTime)
	return
}

func getDirItems(ctx context.Context, tx *sql.Tx, hash string) (dirItems []dao.DirItem, err error) {
	rows, err := tx.QueryContext(ctx, `
		SELECT
			itemHash,
			itemName,
			itemMode,
			itemSize,
			itemCount,
			itemTotalCount,
			itemCreateTime,
			itemModifyTime,
			itemChangeTime,
			itemAccessTime
		FROM _dirItem WHERE Hash=?
	`, hash)
	if err != nil {
		return
	}
	defer rows.Close()
	dirItems = make([]dao.DirItem, 0)
	for rows.Next() {
		var dirItem dao.DirItem
		err = rows.Scan(
			&dirItem.Hash,
			&dirItem.Name,
			&dirItem.Mode,
			&dirItem.Size,
			&dirItem.Count,
			&dirItem.TotalCount,
			&dirItem.CreateTime,
			&dirItem.ModifyTime,
			&dirItem.ChangeTime,
			&dirItem.AccessTime)
		if err != nil {
			return
		}
		dirItems = append(dirItems, dirItem)
	}
	return
}

func UpdateDirItemWithTx(ctx context.Context, tx *sql.Tx, db DbImpl, branchName string, splitPath []string, fn func([][]dao.DirItem) ([]dao.DirItem, error)) (commit dao.Commit, branch dao.Branch, err error) {
	hash, err := getBranchCommitHash(ctx, tx, branchName)
	if err != nil {
		return
	}
	dirItems, err := getDirItems(ctx, tx, hash)
	if err != nil {
		return
	}
	dirItemsList := [][]dao.DirItem{dirItems}
	for i := range splitPath[:len(splitPath)-1] {
		hash, err = getDirItemHash(ctx, tx, hash, splitPath, i)
		if err != nil {
			return
		}
		dirItems, err = getDirItems(ctx, tx, hash)
		if err != nil {
			return
		}
		dirItemsList = append(dirItemsList, dirItems)
	}
	insertDirItems, err := fn(dirItemsList)
	if err != nil {
		return
	}
	i := len(dirItemsList) - 1
	dir, err := db.InsertDirWithTx(ctx, tx, dirItemsList[i], insertDirItems)
	if err != nil {
		return
	}
	for i--; i >= 0; i-- {
		for j := range dirItemsList[i] {
			if dirItemsList[i][j].Name == splitPath[i] {
				dirItemsList[i][j].Hash = dir.Hash()
				dirItemsList[i][j].Size = dir.Size()
				dirItemsList[i][j].Count = dir.Count()
				break
			}
		}
		dir, err = db.InsertDirWithTx(ctx, tx, dirItemsList[i], nil)
		if err != nil {
			return
		}
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

func UpdateDirItemsWithTx(ctx context.Context, tx *sql.Tx, db DbImpl, branchName string, splitPath []string, fn func(*[]dao.DirItem) ([]dao.DirItem, error)) (commit dao.Commit, branch dao.Branch, err error) {
	hash, err := getBranchCommitHash(ctx, tx, branchName)
	if err != nil {
		return
	}
	dirItems, err := getDirItems(ctx, tx, hash)
	if err != nil {
		return
	}
	dirItemsList := [][]dao.DirItem{dirItems}
	for i := range splitPath {
		hash, err = getDirItemHash(ctx, tx, hash, splitPath, i)
		if err != nil {
			return
		}
		dirItems, err = getDirItems(ctx, tx, hash)
		if err != nil {
			return
		}
		dirItemsList = append(dirItemsList, dirItems)
	}
	insertDirItems, err := fn(&dirItems)
	if err != nil {
		return
	}
	i := len(dirItemsList) - 1
	dir, err := db.InsertDirWithTx(ctx, tx, dirItems, insertDirItems)
	if err != nil {
		return
	}
	for i--; i >= 0; i-- {
		for j := range dirItemsList[i] {
			if dirItemsList[i][j].Name == splitPath[i] {
				dirItemsList[i][j].Hash = dir.Hash()
				dirItemsList[i][j].Size = dir.Size()
				dirItemsList[i][j].Count = dir.Count()
				break
			}
		}
		dir, err = db.InsertDirWithTx(ctx, tx, dirItemsList[i], nil)
		if err != nil {
			return
		}
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
