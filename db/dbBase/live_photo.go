package dbBase

import (
	"context"
	"database/sql"
)

func UpsertLivePhoto(ctx context.Context, txOrDb TxOrDb, movHash string, heicHash string, jpgHash string, livpHash string) (err error) {
	if livpHash != "" {
		_, err = txOrDb.ExecContext(ctx, `
		INSERT INTO _live_photo (
			movHash,
			heicHash,
			livpHash
		) VALUES (?, ?, ?) ON CONFLICT DO UPDATE SET
			movHash=?,
			heicHash=?,
			livpHash=?;
		`, movHash, heicHash, livpHash, movHash, heicHash, livpHash)
	} else if heicHash != "" && jpgHash != "" {
		_, err = txOrDb.ExecContext(ctx, `
		INSERT INTO _live_photo (
			movHash,
			heicHash,
			jpgHash
		) VALUES (?, ?, ?) ON CONFLICT DO UPDATE SET
			movHash=?,
			heicHash=?,
			jpgHash=?;
		`, movHash, heicHash, jpgHash, movHash, heicHash, jpgHash)
	} else if heicHash != "" {
		_, err = txOrDb.ExecContext(ctx, `
		INSERT INTO _live_photo (
			movHash,
			heicHash
		) VALUES (?, ?) ON CONFLICT DO UPDATE SET
			movHash=?,
			heicHash=?;
		`, movHash, heicHash, movHash, heicHash)
	} else {
		_, err = txOrDb.ExecContext(ctx, `
		INSERT INTO _live_photo (
			movHash,
			jpgHash
		) VALUES (?, ?) ON CONFLICT DO UPDATE SET
			movHash=?,
			jpgHash=?;
		`, movHash, jpgHash, movHash, jpgHash)
	}
	if err != nil {
		return err
	}
	return err
}

func ListLivePhotoNew(ctx context.Context, conn *sql.DB) (hashList []string, err error) {
	rows, err := conn.QueryContext(ctx, `
	SELECT
		_file_type.hash
	FROM _file_type INNER JOIN _driver_file WHERE _file_type.hash=_driver_file.hash AND _file_type.SubType='zip' AND
		case when _driver_file.name like '%.%' then lower(replace(_driver_file.name, rtrim(_driver_file.name, replace(_driver_file.name, '.', '' ) ), '')) else '' end='livp'
	EXCEPT SELECT livpHash FROM _live_photo
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

func ListLivePhotoAll(ctx context.Context, conn *sql.DB) (hashList []string, err error) {
	rows, err := conn.QueryContext(ctx, `
	SELECT
		_file_type.hash
	FROM _file_type INNER JOIN _driver_file WHERE _file_type.hash=_driver_file.hash AND _file_type.SubType='zip' AND
		case when _driver_file.name like '%.%' then lower(replace(_driver_file.name, rtrim(_driver_file.name, replace(_driver_file.name, '.', '' ) ), '')) else '' end='livp'
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
