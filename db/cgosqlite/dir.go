package cgosqlite

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"strings"

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
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
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
	if len(insertDirItems) == 0 {
		return
	}
	column := 11
	totalRow := len(insertDirItems)
	repeat := 0
	remainRow := totalRow
	maxRow := 32766 / column
	if totalRow > maxRow {
		repeat = totalRow / maxRow
		remainRow = totalRow - repeat*maxRow
		var query string
		query, err = getInsertDirItemQuery(maxRow)
		if err != nil {
			return
		}
		var stmt *sql.Stmt
		stmt, err = tx.PrepareContext(ctx, query)
		if err != nil {
			return
		}
		defer func() {
			if err == nil {
				err = stmt.Close()
			}
		}()
		for i := 0; i < repeat; i++ {
			args := make([]interface{}, maxRow*column)
			for i, dirItem := range insertDirItems[i*maxRow : (i+1)*maxRow] {
				args[i*column] = dir.Hash()
				args[i*column+1] = dirItem.Hash
				args[i*column+2] = dirItem.Name
				args[i*column+3] = dirItem.Mode
				args[i*column+4] = dirItem.Size
				args[i*column+5] = dirItem.Count
				args[i*column+6] = dirItem.TotalCount
				args[i*column+7] = dirItem.CreateTime
				args[i*column+8] = dirItem.ModifyTime
				args[i*column+9] = dirItem.ChangeTime
				args[i*column+10] = dirItem.AccessTime
			}
			// TODO: override if duplicated
			_, err = stmt.ExecContext(ctx, args...)
			if err != nil {
				return
			}
		}
	}
	if remainRow > 0 {
		var query string
		query, err = getInsertDirItemQuery(remainRow)
		if err != nil {
			return
		}
		args := make([]interface{}, remainRow*column)
		for i, dirItem := range insertDirItems[repeat*maxRow:] {
			args[i*column] = dir.Hash()
			args[i*column+1] = dirItem.Hash
			args[i*column+2] = dirItem.Name
			args[i*column+3] = dirItem.Mode
			args[i*column+4] = dirItem.Size
			args[i*column+5] = dirItem.Count
			args[i*column+6] = dirItem.TotalCount
			args[i*column+7] = dirItem.CreateTime
			args[i*column+8] = dirItem.ModifyTime
			args[i*column+9] = dirItem.ChangeTime
			args[i*column+10] = dirItem.AccessTime
		}
		// TODO: override if duplicated
		_, err = tx.ExecContext(ctx, query, args...)
		if err != nil {
			return
		}
	}
	if err != nil {
		return
	}
	return
}

func getInsertDirItemQuery(row int) (string, error) {
	var qs strings.Builder
	_, err := qs.WriteString(`
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
	qs.WriteString(";")
	return qs.String(), err
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

func (db *DB) getDirItem(ctx context.Context, tx *sql.Tx, hash string, splitPath []string, i int) (dirItem dao.DirItem, err error) {
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
