package dbBase

import (
	"context"
	"database/sql"
	"github.com/lazyxu/kfs/dao"
)

func List(ctx context.Context, conn *sql.DB, branchName string, splitPath []string) (dirItems []dao.DirItem, err error) {
	tx, err := conn.Begin()
	if err != nil {
		return
	}
	defer func() {
		err = CommitAndRollback(tx, err)
	}()
	hash, err := getBranchCommitHash(ctx, tx, branchName)
	if err != nil {
		return
	}
	for i := range splitPath {
		hash, err = getDirItemHash(ctx, tx, hash, splitPath, i)
		if err != nil {
			return
		}
	}
	dirItems, err = getDirItems(ctx, tx, hash)
	if err != nil {
		return
	}
	return
}

func ListByHash(ctx context.Context, conn *sql.DB, hash string) (dirItems []dao.DirItem, err error) {
	tx, err := conn.Begin()
	if err != nil {
		return
	}
	defer func() {
		err = CommitAndRollback(tx, err)
	}()
	dirItems, err = getDirItems(ctx, tx, hash)
	if err != nil {
		return
	}
	return
}

func ListV2(ctx context.Context, conn *sql.DB, driverName string, filePath []string) (files []dao.DriverFile, err error) {
	rows, err := conn.QueryContext(ctx, `
		SELECT name,
			hash,
			mode,
			size,
			createTime,
			modifyTime,
			changeTime,
			accessTime
		FROM _driver_file WHERE driver_name=? and dirpath=?
	`, driverName, arrayToJson(filePath))
	if err != nil {
		return
	}
	defer rows.Close()
	files = make([]dao.DriverFile, 0)
	for rows.Next() {
		var file dao.DriverFile
		err = rows.Scan(
			&file.Name,
			&file.Hash,
			&file.Mode,
			&file.Size,
			&file.CreateTime,
			&file.ModifyTime,
			&file.ChangeTime,
			&file.AccessTime)
		if err != nil {
			return
		}
		files = append(files, file)
	}
	return
}
