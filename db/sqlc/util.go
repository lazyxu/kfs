package sqlc

import (
	"database/sql"
	"time"
)

func NInt64(v int64) sql.NullInt64 {
	return sql.NullInt64{Int64: v, Valid: true}
}

func NString(v string) sql.NullString {
	return sql.NullString{String: v, Valid: true}
}

func NFloat64(v float64) sql.NullFloat64 {
	return sql.NullFloat64{Float64: v, Valid: true}
}

func NTime(v time.Time) sql.NullTime {
	return sql.NullTime{Time: v, Valid: true}
}
