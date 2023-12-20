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

func GetEarliestCrated(ctx context.Context, conn *sql.DB, db DbImpl, hash string) int64 {
	return 0
}
