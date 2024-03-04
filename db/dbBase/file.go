package dbBase

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
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

var ErrNoSuchFileOrDir = errors.New("no such file or dir")

func GetDriverFile(ctx context.Context, conn *sql.DB, driverId uint64, filePath []string) (file dao.DriverFile, err error) {
	if len(filePath) == 0 {
		err = errors.New("/: Is a directory")
		return
	}
	file.DriverId = driverId
	file.DirPath = filePath[:len(filePath)-1]
	file.Name = filePath[len(filePath)-1]

	rows, err := conn.QueryContext(ctx, `
		SELECT hash,
		mode,
		size,
		createTime,
		modifyTime,
		changeTime,
		accessTime
		FROM _driver_file WHERE driverId=? and dirPath=? and name=?
	`, file.DriverId, arrayToJson(file.DirPath), file.Name)
	if err != nil {
		return
	}
	defer rows.Close()
	if !rows.Next() {
		err = ErrNoSuchFileOrDir
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

func jsonToArray(data []byte) (arr []string) {
	err := json.Unmarshal(data, &arr)
	if err != nil {
		panic(err)
	}
	return arr
}

func UpsertDriverFile(ctx context.Context, conn *sql.DB, f dao.DriverFile) error {
	tx, err := conn.Begin()
	if err != nil {
		return err
	}
	defer func() {
		err = CommitAndRollback(tx, err)
	}()
	res, err := tx.ExecContext(ctx, `
	INSERT INTO _driver_file (
		driverId,
		dirPath,
		name,
		hash,
		mode,
		size,
		createTime,
		modifyTime,
		changeTime,
		accessTime
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?) ON CONFLICT DO UPDATE SET
		hash=?,
		mode=?,
		size=?,
		createTime=?,
		modifyTime=?,
		changeTime=?,
		accessTime=?
	WHERE
		hash!=? OR
		mode!=? OR
		size!=? OR
		createTime!=? OR
		modifyTime!=? OR
		changeTime!=? OR
		accessTime!=?;
	`, f.DriverId, arrayToJson(f.DirPath), f.Name, f.Hash, f.Mode, f.Size, f.CreateTime, f.ModifyTime, f.ChangeTime, f.AccessTime,
		f.Hash, f.Mode, f.Size, f.CreateTime, f.ModifyTime, f.ChangeTime, f.AccessTime,
		f.Hash, f.Mode, f.Size, f.CreateTime, f.ModifyTime, f.ChangeTime, f.AccessTime)
	if err != nil {
		return err
	}
	i, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if i > 0 {
		fmt.Printf("%+v %s: %d\n", f.DirPath, f.Name, i)
		_, err = tx.ExecContext(ctx, `
	INSERT INTO _driver_file_history (
		driverId,
		dirPath,
		name,
		hash,
		mode,
		size,
		createTime,
		modifyTime,
		changeTime,
		accessTime,
	    uploadDeviceId,
	    uploadTime
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);
	`, f.DriverId, arrayToJson(f.DirPath), f.Name, f.Hash, f.Mode, f.Size, f.CreateTime, f.ModifyTime, f.ChangeTime, f.AccessTime, f.UploadDeviceId, f.UploadTime)
		if err != nil {
			return err
		}
	}
	return nil
}

func upsertDriverFilesQuery(row int) (string, error) {
	var qs strings.Builder
	_, err := qs.WriteString(`
	INSERT OR REPLACE INTO _driver_file (
		driverId,
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
	) VALUES `)
	if err != nil {
		return "", err
	}
	for i := 0; i < row; i++ {
		if i != 0 {
			qs.WriteString(", ")
		}
		qs.WriteString("(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	}
	qs.WriteString(`;`)
	return qs.String(), err
}

func UpsertDriverFiles(ctx context.Context, conn *sql.DB, db DbImpl, files []dao.DriverFile) error {
	tx, err := conn.Begin()
	if err != nil {
		return err
	}
	defer func() {
		err = CommitAndRollback(tx, err)
	}()
	err = InsertBatch[dao.DriverFile](ctx, tx, db.MaxBatchSize(), files, 11, upsertDriverFilesQuery, func(args []interface{}, start int, f dao.DriverFile) {
		args[start] = f.DriverId
		args[start+1] = arrayToJson(f.DirPath)
		args[start+2] = f.Name
		args[start+3] = f.Version
		args[start+4] = f.Hash
		args[start+5] = f.Mode
		args[start+6] = f.Size
		args[start+7] = f.CreateTime
		args[start+8] = f.ModifyTime
		args[start+9] = f.ChangeTime
		args[start+10] = f.AccessTime
	})
	if err != nil {
		return err
	}
	return err
}

func UpsertDriverFileMysql(ctx context.Context, txOrDb TxOrDb, f dao.DriverFile) error {
	_, err := txOrDb.ExecContext(ctx, `
	INSERT INTO _driver_file (
		driverId,
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
	`, f.DriverId, f.DirPath, f.Name, f.Version, f.Hash, f.Mode, f.Size, f.CreateTime, f.ModifyTime, f.ChangeTime, f.AccessTime,
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

func InsertFileMd5(ctx context.Context, conn *sql.DB, db DbImpl, hash string, hashMd5 string) error {
	// TODO: on duplicated key check size.
	_, err := conn.ExecContext(ctx, `
	INSERT INTO _file_md5 (
		hash,
		md5
	) VALUES (?, ?);
	`, hash, hashMd5)
	if err != nil {
		if db.IsUniqueConstraintError(err) {
			return nil
		}
		return err
	}
	return err
}

func getListFileMd5Query(row int) (string, error) {
	var qs strings.Builder
	_, err := qs.WriteString(`
	SELECT
		hash,
		md5
	FROM _file_md5 WHERE md5 IN (`)
	if err != nil {
		return "", err
	}
	for i := 0; i < row; i++ {
		if i != 0 {
			qs.WriteString(", ")
		}
		qs.WriteString("?")
	}
	qs.WriteString(");")
	return qs.String(), err
}

func toAny(s []string) []any {
	c := make([]any, len(s))
	for i, v := range s {
		c[i] = v
	}
	return c
}

func ListFileMd5(ctx context.Context, conn *sql.DB, md5List []string) (m map[string]string, err error) {
	query, err := getListFileMd5Query(len(md5List))
	if err != nil {
		return
	}
	rows, err := conn.QueryContext(ctx, query, toAny(md5List)...)
	if err != nil {
		return
	}
	defer rows.Close()
	m = make(map[string]string)
	for rows.Next() {
		var hash, hashMd5 string
		err = rows.Scan(
			&hash,
			&hashMd5)
		if err != nil {
			return
		}
		m[hashMd5] = hash
	}
	return
}

func SumFileSize(ctx context.Context, conn *sql.DB) (size uint64, err error) {
	// TODO: on duplicated key check size.
	rows, err := conn.QueryContext(ctx, `
	SELECT SUM(size) FROM _file;
	`)
	if err != nil {
		return
	}
	defer rows.Close()
	if rows.Next() {
		err = rows.Scan(&size)
		if err != nil {
			return
		}
	} else {
		err = ErrNoRecords
	}
	return
}

func ListFile(ctx context.Context, conn *sql.DB) (hashList []string, err error) {
	rows, err := conn.QueryContext(ctx, `
	SELECT hash FROM _file;
	`)
	if err != nil {
		return
	}
	defer rows.Close()
	hashList = []string{}
	for rows.Next() {
		var hash string
		err = rows.Scan(&hash)
		if err != nil {
			return
		}
		hashList = append(hashList, hash)
	}
	return
}
