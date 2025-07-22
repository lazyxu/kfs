package sqlc

import (
	"database/sql"
	"time"

	"github.com/lazyxu/kfs/dao"
)

func NInt64(v int64) dao.NullInt64 {
	return dao.NullInt64{NullInt64: sql.NullInt64{Int64: v, Valid: true}}
}

func NInt64Null() dao.NullInt64 {
	return dao.NullInt64{NullInt64: sql.NullInt64{Valid: false}}
}

func NBoolNull() dao.NullBool {
	return dao.NullBool{NullBool: sql.NullBool{Valid: false}}
}

func NString(v string) dao.NullString {
	return dao.NullString{NullString: sql.NullString{String: v, Valid: true}}
}

func NFloat64(v float64) dao.NullFloat64 {
	return dao.NullFloat64{NullFloat64: sql.NullFloat64{Float64: v, Valid: true}}
}

func NTime(v time.Time) dao.NullTime {
	return dao.NullTime{NullTime: sql.NullTime{Time: v, Valid: true}}
}
