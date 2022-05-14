package noncgo

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"io"
	"os"
	"strings"
)

type Dir struct {
	fileOrDir
	count      uint64
	totalCount uint64
}

func (dir Dir) Count() uint64 {
	return dir.count
}

func (dir Dir) TotalCount() uint64 {
	return dir.totalCount
}

func NewDir(hash string, size uint64, count uint64, totalCount uint64) Dir {
	return Dir{fileOrDir{hash, size}, count, totalCount}
}

// https://zhuanlan.zhihu.com/p/343682839
type DirItem struct {
	Hash       string
	Name       string
	Mode       uint64
	Size       uint64
	Count      uint64
	TotalCount uint64
	CreateTime uint64 // linux does not support it.
	ModifyTime uint64
	ChangeTime uint64 // windows does not support it.
	AccessTime uint64
}

func NewDirItem(fileOrDir FileOrDir, name string, mode uint64, createTime uint64, modifyTime uint64, changeTime uint64, accessTime uint64) DirItem {
	return DirItem{fileOrDir.Hash(), name, mode, fileOrDir.Size(), fileOrDir.Count(), fileOrDir.TotalCount(), createTime, modifyTime, changeTime, accessTime}
}

func writeMutil(w io.Writer, order binary.ByteOrder, data []any) {
	for _, v := range data {
		err := binary.Write(w, order, v)
		if err != nil {
			panic(err)
		}
	}
}

func (dir *Dir) Cal(dirItems []DirItem) {
	hash := sha256.New()
	err := binary.Write(hash, binary.LittleEndian, uint64(len(dirItems)))
	if err != nil {
		panic(err)
	}
	dir.size = 0
	dir.count = uint64(len(dirItems))
	for _, dirItem := range dirItems {
		writeMutil(hash, binary.LittleEndian, []any{
			[]byte(dirItem.Hash),
			[]byte(dirItem.Name),
			dirItem.Mode,
			dirItem.CreateTime,
			dirItem.ModifyTime,
			dirItem.ChangeTime,
			dirItem.AccessTime,
		})
		dir.size += dirItem.Size
		dir.totalCount += dirItem.TotalCount
	}
	dir.hash = hex.EncodeToString(hash.Sum(nil))
}

func (db *DB) WriteDir(ctx context.Context, dirItems []DirItem) (dir Dir, err error) {
	tx, err := db._db.Begin()
	if err != nil {
		return
	}
	defer func() {
		if err != nil {
			println(err.Error())
			err = tx.Rollback()
			if err != sql.ErrTxDone {
				return
			}
		}
	}()
	defer func() {
		if err == nil {
			err = tx.Commit()
		}
	}()
	return db.writeDir(ctx, tx, dirItems)
}

type TxOrDb interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
}

func (db *DB) writeDir(ctx context.Context, tx TxOrDb, dirItems []DirItem) (dir Dir, err error) {
	dir.Cal(dirItems)
	// TODO: error if size or count is not equal
	_, err = tx.ExecContext(ctx, `
	INSERT INTO dir VALUES (?, ?, ?, ?);
	`, dir.hash, dir.size, dir.count, dir.totalCount)
	if err != nil {
		if isUniqueConstraintError(err) {
			err = nil
		}
		return
	}
	stmt, err := tx.PrepareContext(ctx, `
	INSERT INTO dirItem (
		hash,
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
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);
	`)
	if err != nil {
		return
	}
	defer func() {
		if err == nil {
			err = stmt.Close()
		}
	}()
	for _, dirItem := range dirItems {
		// TODO: override if duplicated
		_, err = stmt.ExecContext(ctx,
			dir.hash,
			dirItem.Hash,
			dirItem.Name,
			dirItem.Mode,
			dirItem.Size,
			dirItem.Count,
			dirItem.TotalCount,
			dirItem.CreateTime,
			dirItem.ModifyTime,
			dirItem.ChangeTime,
			dirItem.AccessTime)
		if err != nil {
			return
		}
	}
	return
}

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

func (db *DB) GetFileHash(ctx context.Context, branchName string, splitPath []string) (hash string, err error) {
	tx, err := db._db.Begin()
	if err != nil {
		return
	}
	defer func() {
		if err != nil {
			println(err.Error())
			err = tx.Rollback()
			if err != sql.ErrTxDone {
				return
			}
		}
	}()
	defer func() {
		if err == nil {
			err = tx.Commit()
		}
	}()
	hash, err = db.getBranchCommitHash(ctx, tx, branchName)
	if err != nil {
		return
	}
	for i := range splitPath[:len(splitPath)-1] {
		hash, err = db.getDirItemHash(ctx, tx, hash, splitPath, i)
		if err != nil {
			return
		}
	}
	hash, mode, err := db.getDirItemHashMode(ctx, tx, hash, splitPath, len(splitPath)-1)
	if err != nil {
		return
	}
	if os.FileMode(mode).IsDir() {
		return "", errors.New("/" + strings.Join(splitPath, "/") + ": Is a directory")
	}
	return
}

func (db *DB) getBranchCommitHash(ctx context.Context, tx *sql.Tx, branchName string) (hash string, err error) {
	rows, err := tx.QueryContext(ctx, `
		SELECT [commit].hash FROM branch INNER JOIN [commit] WHERE branch.name=? and [commit].id=branch.commitId
	`, branchName)
	if err != nil {
		return
	}
	defer rows.Close()
	if !rows.Next() {
		return "", errors.New("no such branch " + branchName)
	}
	err = rows.Scan(&hash)
	return
}

func (db *DB) getDirItemHash(ctx context.Context, tx *sql.Tx, hash string, splitPath []string, i int) (itemHash string, err error) {
	rows, err := tx.QueryContext(ctx, `
		SELECT itemHash FROM dirItem WHERE hash=? and itemName=?
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

func (db *DB) getDirItemHashMode(ctx context.Context, tx *sql.Tx, hash string, splitPath []string, i int) (itemHash string, itemMode uint64, err error) {
	rows, err := tx.QueryContext(ctx, `
		SELECT itemHash, itemMode  FROM dirItem WHERE hash=? and itemName=?
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

func (db *DB) getDirItems(ctx context.Context, tx *sql.Tx, hash string) (dirItems []DirItem, err error) {
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
		FROM dirItem WHERE hash=?
	`, hash)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var dirItem DirItem
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

func (db *DB) Remove(ctx context.Context, branchName string, splitPath []string) (err error) {
	tx, err := db._db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			println(err.Error())
			err = tx.Rollback()
			if err != sql.ErrTxDone {
				return
			}
		}
	}()
	defer func() {
		if err == nil {
			err = tx.Commit()
		}
	}()
	hash, err := db.getBranchCommitHash(ctx, tx, branchName)
	if err != nil {
		return
	}
	if len(splitPath) == 0 {
		var dir Dir
		dir, err = db.WriteDir(ctx, nil)
		if err != nil {
			return
		}
		commit := NewCommit(dir, branchName)
		err = db.writeCommit(ctx, tx, &commit)
		if err != nil {
			return
		}
		branch := NewBranch(branchName, "", commit, dir)
		err = db.writeBranch(ctx, tx, branch)
		return err
	}
	dirItems, err := db.getDirItems(ctx, tx, hash)
	if err != nil {
		return
	}
	dirItemsList := [][]DirItem{dirItems}
	for i := range splitPath[:len(splitPath)-1] {
		hash, err = db.getDirItemHash(ctx, tx, hash, splitPath, i)
		if err != nil {
			return
		}
		dirItems, err = db.getDirItems(ctx, tx, hash)
		if err != nil {
			return
		}
		dirItemsList = append(dirItemsList, dirItems)
	}
	i := len(dirItemsList) - 1
	find := false
	for j, dirItem := range dirItemsList[i] {
		if dirItem.Name == splitPath[i] {
			dirItemsList[i][j] = dirItemsList[i][len(dirItemsList[i])-1]
			dirItemsList[i] = dirItemsList[i][:len(dirItemsList[i])-1]
			find = true
			break
		}
	}
	if !find {
		return errors.New("no such file or dir: /" + strings.Join(splitPath, "/"))
	}
	dir, err := db.writeDir(ctx, tx, dirItemsList[i])
	if err != nil {
		return
	}
	for i--; i >= 0; i-- {
		for j := range dirItemsList[i] {
			if dirItemsList[i][j].Name == splitPath[i] {
				dirItemsList[i][j].Hash = dir.hash
				dirItemsList[i][j].Size = dir.size
				dirItemsList[i][j].Count = dir.count
				break
			}
		}
		dir, err = db.writeDir(ctx, tx, dirItemsList[i])
		if err != nil {
			return
		}
	}
	commit := NewCommit(dir, branchName)
	err = db.writeCommit(ctx, tx, &commit)
	if err != nil {
		return
	}
	branch := NewBranch(branchName, "", commit, dir)
	err = db.writeBranch(ctx, tx, branch)
	return
}
