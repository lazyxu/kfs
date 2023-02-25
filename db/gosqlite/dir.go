package gosqlite

import (
	"context"
	"database/sql"
	"errors"
	"github.com/lazyxu/kfs/db/dbBase"
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
		err = CommitAndRollback(tx, err)
	}()
	return db.InsertDirWithTx(ctx, tx, dirItems, dirItems)
}

func (db *DB) InsertDirWithTx(ctx context.Context, tx *sql.Tx, dirItems []dao.DirItem, insertDirItems []dao.DirItem) (dir dao.Dir, err error) {
	dir.Cal(dirItems)
	// TODO: error if size or count is not equal
	_, err = tx.ExecContext(ctx, `
	INSERT INTO _dir VALUES (?, ?, ?, ?);
	`, dir.Hash(), dir.Size(), dir.Count(), dir.TotalCount())
	if err != nil {
		if db.IsUniqueConstraintError(err) {
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

func (db *DB) RemoveDirItem(ctx context.Context, branchName string, splitPath []string) (commit dao.Commit, branch dao.Branch, err error) {
	conn := db.getConn()
	defer db.putConn(conn)
	tx, err := conn.Begin()
	if err != nil {
		return
	}
	defer func() {
		err = CommitAndRollback(tx, err)
	}()
	return db.UpdateDirItemWithTx(ctx, tx, branchName, splitPath, func(dirItemsList [][]dao.DirItem) ([]dao.DirItem, error) {
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

func (db *DB) UpdateDirItemWithTx(ctx context.Context, tx *sql.Tx, branchName string, splitPath []string, fn func([][]dao.DirItem) ([]dao.DirItem, error)) (commit dao.Commit, branch dao.Branch, err error) {
	return dbBase.UpdateDirItemWithTx(ctx, tx, db, branchName, splitPath, fn)
}

func (db *DB) UpdateDirItemsWithTx(ctx context.Context, tx *sql.Tx, branchName string, splitPath []string, fn func(*[]dao.DirItem) ([]dao.DirItem, error)) (commit dao.Commit, branch dao.Branch, err error) {
	return dbBase.UpdateDirItemsWithTx(ctx, tx, db, branchName, splitPath, fn)
}
