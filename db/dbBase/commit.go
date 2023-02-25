package dbBase

import (
	"context"
	"database/sql"
	"errors"
	"github.com/lazyxu/kfs/dao"
)

func InsertCommitWithTxOrDb(ctx context.Context, txOrDb TxOrDb, commit *dao.Commit) error {
	// TODO: if Hash not changed.
	res, err := txOrDb.ExecContext(ctx, `
	INSERT INTO _commit (createTime, Hash, lastId)
	VALUES (?, ?, ifnull((SELECT commitId FROM _branch WHERE _branch.name=?), 0));;
	`, commit.CreateTime(), commit.Hash, commit.BranchName())
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	commit.Id = uint64(id)
	return err
}

func InsertCommitWithTxOrDbCgoSqlite(ctx context.Context, txOrDb TxOrDb, commit *dao.Commit) error {
	// TODO: if Hash not changed.
	res, err := txOrDb.ExecContext(ctx, `
	INSERT INTO _commit (createTime, Hash, lastId)
	VALUES (?, ?, ifnull((SELECT commitId FROM _branch WHERE _branch.name=?), 0));;
	`, commit.CreateTime(), commit.Hash, commit.BranchName())
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		rows, err := txOrDb.QueryContext(ctx, `
	SELECT id FROM _commit ORDER BY id DESC LIMIT 1
	`, commit.CreateTime(), commit.Hash, commit.BranchName())
		if err != nil {
			return err
		}
		defer rows.Close()
		if !rows.Next() {
			return errors.New("no commit")
		}
		err = rows.Scan(&commit.Id)
		if err != nil {
			return err
		}
	} else {
		id, err := res.LastInsertId()
		if err != nil {
			return err
		}
		commit.Id = uint64(id)
	}
	return err
}

func getBranchCommitHash(ctx context.Context, tx *sql.Tx, branchName string) (hash string, err error) {
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
