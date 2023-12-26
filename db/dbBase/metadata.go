package dbBase

import (
	"context"
	"database/sql"
	"time"
)

func InsertDCIMMetadataTime(ctx context.Context, conn *sql.DB, db DbImpl, hash string, t int64) (exist bool, err error) {
	tt := time.Unix(0, t)
	_, err = conn.ExecContext(ctx, `
	INSERT INTO _dcim_metadata_time (
		hash,
		time,
		year,
		month,
		day
	) VALUES (?, ?, ?, ?, ?)`, hash, t, tt.Year(), tt.Month(), tt.Day())
	if db.IsUniqueConstraintError(err) {
		exist = true
		err = nil
	}
	return
}

func UpsertDCIMMetadataTime(ctx context.Context, conn *sql.DB, hash string, t int64) error {
	tt := time.Unix(0, t)
	_, err := conn.ExecContext(ctx, `
	INSERT INTO _dcim_metadata_time (
		hash,
		time,
		year,
		month,
		day
	) VALUES (?, ?, ?, ?, ?) ON CONFLICT(hash) DO UPDATE SET
		time=?,
		year=?,
		month=?,
		day=?;`, hash, t, tt.Year(), tt.Month(), tt.Day(), t, tt.Year(), tt.Month(), tt.Day())
	if err != nil {
		return err
	}
	return nil
}

func GetEarliestCrated(ctx context.Context, conn *sql.DB, db DbImpl, hash string) (t int64, err error) {
	rows, err := conn.QueryContext(ctx, `
		SELECT min(min(createTime, modifyTime, accessTime)) FROM _driver_file WHERE hash=?;
	`, hash)
	if err != nil {
		return
	}
	defer rows.Close()
	if rows.Next() {
		err = rows.Scan(&t)
	}
	return
}
