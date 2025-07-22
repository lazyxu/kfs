package dao

import (
	"database/sql"
	"encoding/json"
)

type NullInt64 struct {
	sql.NullInt64
}

func (n NullInt64) MarshalJSON() ([]byte, error) {
	if !n.Valid {
		return []byte("null"), nil // 返回 null 或完全跳过字段（需配合 omitempty）
	}
	return json.Marshal(n.Int64) // 有效值序列化为数字
}

type NullBool struct {
	sql.NullBool
}

func (n NullBool) MarshalJSON() ([]byte, error) {
	if !n.Valid {
		return []byte("null"), nil // 返回 null 或完全跳过字段（需配合 omitempty）
	}
	return json.Marshal(n.Bool) // 有效值序列化为数字
}

type NullString struct {
	sql.NullString
}

func (n NullString) MarshalJSON() ([]byte, error) {
	if !n.Valid {
		return []byte("null"), nil // 返回 null 或完全跳过字段（需配合 omitempty）
	}
	return json.Marshal(n.String) // 有效值序列化为数字
}

type NullFloat64 struct {
	sql.NullFloat64
}

func (n NullFloat64) MarshalJSON() ([]byte, error) {
	if !n.Valid {
		return []byte("null"), nil // 返回 null 或完全跳过字段（需配合 omitempty）
	}
	return json.Marshal(n.Float64) // 有效值序列化为数字
}

type NullTime struct {
	sql.NullTime
}

func (n NullTime) MarshalJSON() ([]byte, error) {
	if !n.Valid {
		return []byte("null"), nil // 返回 null 或完全跳过字段（需配合 omitempty）
	}
	return json.Marshal(n.Time) // 有效值序列化为数字
}
