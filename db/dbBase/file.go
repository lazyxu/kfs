package dbBase

import (
	"context"
	"database/sql"
	"encoding/json"
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

func GetDriverFile(ctx context.Context, conn *sql.DB, driverName string, splitPath []string) (file dao.DriverFile, err error) {
	if len(splitPath) == 0 {
		err = errors.New("/: Is a directory")
		return
	}
	file.DriverName = driverName
	file.DirPath = splitPath[:len(splitPath)-1]
	file.Name = splitPath[len(splitPath)-1]
	file.Version = 0

	rows, err := conn.QueryContext(ctx, `
		SELECT hash,
		mode,
		size,
		createTime,
		modifyTime,
		changeTime,
		accessTime
		FROM _driver_file WHERE driverName=? and dirPath=? and name=? and version=0
	`, file.DriverName, arrayToJson(file.DirPath), file.Name)
	if err != nil {
		return
	}
	defer rows.Close()
	if !rows.Next() {
		err = errors.New("no such file or dir: /" + strings.Join(splitPath, "/"))
		return
	}
	err = rows.Scan(
		&file.Hash,
		&file.Mode,
		&file.Size,
		&file.CreateTime,
		&file.ModifyTime,
		&file.ChangeTime,
		&file.AccessTime)
	return
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

func arrayToJson(arr []string) []byte {
	if arr == nil {
		arr = []string{}
	}
	data, err := json.Marshal(arr)
	if err != nil {
		panic(err)
	}
	return data
}

func UpsertDriverFile(ctx context.Context, txOrDb TxOrDb, f dao.DriverFile) error {
	_, err := txOrDb.ExecContext(ctx, `
	INSERT INTO _driver_file (
		driverName,
		dirPath,
		name,
	    version,
		hash,
		mode,
		size,
		createTime,
		modifyTime,
		changeTime,
		accessTime
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?) ON CONFLICT DO UPDATE SET
		hash=?,
		mode=?,
		size=?,
		createTime=?,
		modifyTime=?,
		changeTime=?,
		accessTime=?;
	`, f.DriverName, arrayToJson(f.DirPath), f.Name, f.Version, f.Hash, f.Mode, f.Size, f.CreateTime, f.ModifyTime, f.ChangeTime, f.AccessTime,
		f.Hash, f.Mode, f.Size, f.CreateTime, f.ModifyTime, f.ChangeTime, f.AccessTime)
	if err != nil {
		return err
	}
	return err
}

func UpsertDriverFileMysql(ctx context.Context, txOrDb TxOrDb, f dao.DriverFile) error {
	_, err := txOrDb.ExecContext(ctx, `
	INSERT INTO _driver_file (
		driverName,
		dirPath,
		name,
	    version,
		hash,
		mode,
		size,
		createTime,
		modifyTime,
		changeTime,
		accessTime
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE 
		hash=?,
		mode=?,
		size=?,
		createTime=?,
		modifyTime=?,
		changeTime=?,
		accessTime=?;
	`, f.DriverName, f.DirPath, f.Name, f.Version, f.Hash, f.Mode, f.Size, f.CreateTime, f.ModifyTime, f.ChangeTime, f.AccessTime,
		f.Hash, f.Mode, f.Size, f.CreateTime, f.ModifyTime, f.ChangeTime, f.AccessTime)
	if err != nil {
		return err
	}
	return err
}

func InsertFile(ctx context.Context, conn *sql.DB, db DbImpl, hash string, size uint64) error {
	// TODO: on duplicated key check size.
	_, err := conn.ExecContext(ctx, `
	INSERT INTO _file VALUES (?, ?);
	`, hash, size)
	if err != nil {
		if db.IsUniqueConstraintError(err) {
			return nil
		}
		return err
	}
	return err
}
