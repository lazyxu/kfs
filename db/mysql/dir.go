package mysql

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"strings"
	"time"

	"github.com/lazyxu/kfs/dao"
)

func (db *DB) WriteDir(ctx context.Context, dirItems []dao.DirItem) (dir dao.Dir, err error) {
	conn := db.getConn()
	defer db.putConn(conn)
	tx, err := conn.Begin()
	if err != nil {
		return
	}
	defer func() {
		err = commitAndRollback(tx, err)
	}()
	return db.writeDir(ctx, tx, dirItems, dirItems)
}

type TxOrDb interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
}

func (db *DB) writeDir(ctx context.Context, tx TxOrDb, dirItems []dao.DirItem, insertDirItems []dao.DirItem) (dir dao.Dir, err error) {
	dir.Cal(dirItems)
	// TODO: error if size or count is not equal
	_, err = tx.ExecContext(ctx, `
	INSERT INTO _dir VALUES (?, ?, ?, ?);
	`, dir.Hash(), dir.Size(), dir.Count(), dir.TotalCount())
	if err != nil {
		if isUniqueConstraintError(err) {
			err = nil
		}
		return
	}
	stmt, err := tx.PrepareContext(ctx, `
	INSERT INTO _dirItem (
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
	for _, dirItem := range insertDirItems {
		// TODO: override if duplicated
		_, err = stmt.ExecContext(ctx,
			dir.Hash(),
			dirItem.Hash,
			dirItem.Name,
			dirItem.Mode,
			dirItem.Size,
			dirItem.Count,
			dirItem.TotalCount,
			time.Unix(0, int64(dirItem.CreateTime)),
			time.Unix(0, int64(dirItem.ModifyTime)),
			time.Unix(0, int64(dirItem.ChangeTime)),
			time.Unix(0, int64(dirItem.AccessTime)))
		if err != nil {
			return
		}
	}
	return
}

func (db *DB) GetFileHash(ctx context.Context, branchName string, splitPath []string) (hash string, err error) {
	conn := db.getConn()
	defer db.putConn(conn)
	tx, err := conn.Begin()
	if err != nil {
		return
	}
	defer func() {
		err = commitAndRollback(tx, err)
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
		SELECT _commit.Hash FROM _branch INNER JOIN _commit WHERE _branch.name=? and _commit.id=_branch.commitId
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

func (db *DB) getDirItemHashMode(ctx context.Context, tx *sql.Tx, hash string, splitPath []string, i int) (itemHash string, itemMode uint64, err error) {
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

func (db *DB) getDirItems(ctx context.Context, tx *sql.Tx, hash string) (dirItems []dao.DirItem, err error) {
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
	for rows.Next() {
		var dirItem dao.DirItem
		var createTime time.Time
		var modifyTime time.Time
		var changeTime time.Time
		var accessTime time.Time
		err = rows.Scan(
			&dirItem.Hash,
			&dirItem.Name,
			&dirItem.Mode,
			&dirItem.Size,
			&dirItem.Count,
			&dirItem.TotalCount,
			&createTime,
			&modifyTime,
			&changeTime,
			&accessTime)
		if err != nil {
			return
		}
		dirItem.CreateTime = uint64(createTime.UnixNano())
		dirItem.ModifyTime = uint64(modifyTime.UnixNano())
		dirItem.ChangeTime = uint64(changeTime.UnixNano())
		dirItem.AccessTime = uint64(accessTime.UnixNano())
		dirItems = append(dirItems, dirItem)
	}
	return
}

func (db *DB) RemoveDirItem(ctx context.Context, branchName string, splitPath []string) (commit dao.Commit, branch dao.Branch, err error) {
	conn := db.getConn()
	defer db.putConn(conn)
	tx, err := conn.Begin()
	if err != nil {
		return
	}
	defer func() {
		err = commitAndRollback(tx, err)
	}()
	return db.updateDirItem(ctx, tx, branchName, splitPath, func(dirItemsList [][]dao.DirItem) ([]dao.DirItem, error) {
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
			return nil, errors.New("no such file or dir: /" + strings.Join(splitPath, "/"))
		}
		return nil, nil
	})
}

func (db *DB) updateDirItem(ctx context.Context, tx *sql.Tx, branchName string, splitPath []string, fn func([][]dao.DirItem) ([]dao.DirItem, error)) (commit dao.Commit, branch dao.Branch, err error) {
	hash, err := db.getBranchCommitHash(ctx, tx, branchName)
	if err != nil {
		return
	}
	dirItems, err := db.getDirItems(ctx, tx, hash)
	if err != nil {
		return
	}
	dirItemsList := [][]dao.DirItem{dirItems}
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
	insertDirItems, err := fn(dirItemsList)
	if err != nil {
		return
	}
	i := len(dirItemsList) - 1
	dir, err := db.writeDir(ctx, tx, dirItemsList[i], insertDirItems)
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
		dir, err = db.writeDir(ctx, tx, dirItemsList[i], nil)
		if err != nil {
			return
		}
	}
	commit = dao.NewCommit(dir, branchName, "")
	err = db.writeCommit(ctx, tx, &commit)
	if err != nil {
		return
	}
	branch = dao.NewBranch(branchName, commit, dir)
	err = db.writeBranch(ctx, tx, branch)
	if err != nil {
		return
	}
	return
}

func (db *DB) updateDirItems(ctx context.Context, tx *sql.Tx, branchName string, splitPath []string, fn func(*[]dao.DirItem) ([]dao.DirItem, error)) (commit dao.Commit, branch dao.Branch, err error) {
	hash, err := db.getBranchCommitHash(ctx, tx, branchName)
	if err != nil {
		return
	}
	dirItems, err := db.getDirItems(ctx, tx, hash)
	if err != nil {
		return
	}
	dirItemsList := [][]dao.DirItem{dirItems}
	for i := range splitPath {
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
	insertDirItems, err := fn(&dirItems)
	if err != nil {
		return
	}
	i := len(dirItemsList) - 1
	dir, err := db.writeDir(ctx, tx, dirItems, insertDirItems)
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
		dir, err = db.writeDir(ctx, tx, dirItemsList[i], nil)
		if err != nil {
			return
		}
	}
	commit = dao.NewCommit(dir, branchName, "")
	err = db.writeCommit(ctx, tx, &commit)
	if err != nil {
		return
	}
	branch = dao.NewBranch(branchName, commit, dir)
	err = db.writeBranch(ctx, tx, branch)
	if err != nil {
		return
	}
	return
}