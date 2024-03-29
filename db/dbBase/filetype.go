package dbBase

import (
	"context"
	"database/sql"
	"github.com/lazyxu/kfs/dao"
)

func InsertFileType(ctx context.Context, conn *sql.DB, db DbImpl, hash string, t dao.FileType) (exist bool, err error) {
	_, err = conn.ExecContext(ctx, `
	INSERT INTO _file_type (
		hash,
		Type,
		SubType,
		Extension
	) VALUES (?, ?, ?, ?)`, hash, t.Type, t.SubType, t.Extension)
	if db.IsUniqueConstraintError(err) {
		exist = true
		err = nil
	}
	return
}

func UpsertFileType(ctx context.Context, conn *sql.DB, hash string, t dao.FileType) error {
	_, err := conn.ExecContext(ctx, `
	INSERT INTO _file_type (
		hash,
		Type,
		SubType,
		Extension
	) VALUES (?, ?, ?, ?) ON CONFLICT(hash) DO UPDATE SET
	    Type=?,
	    SubType=?,
	    Extension=?`, hash, t.Type, t.SubType, t.Extension, t.Type, t.SubType, t.Extension)
	if err != nil {
		return err
	}
	return nil
}

func ListExpectFileType(ctx context.Context, conn *sql.DB) (hashList []string, err error) {
	rows, err := conn.QueryContext(ctx, `
	SELECT hash FROM _file EXCEPT SELECT hash FROM _file_type;
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

func ListFileHash(ctx context.Context, conn *sql.DB) (hashList []string, err error) {
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

func GetFileType(ctx context.Context, conn *sql.DB, hash string) (fileType dao.FileType, err error) {
	rows, err := conn.QueryContext(ctx, `
	SELECT 
		Type,
		SubType,
		Extension
	FROM _file_type WHERE hash=?;
	`, hash)
	if err != nil {
		return
	}
	defer rows.Close()
	if rows.Next() {
		err = rows.Scan(&fileType.Type, &fileType.SubType, &fileType.Extension)
		if err != nil {
			return
		}
	} else {
		err = ErrNoRecords
	}
	return
}
