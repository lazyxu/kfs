package sqlc

import (
	"context"
	"fmt"
	"strings"
)

func getListFileMd5Query(row int) (string, error) {
	var qs strings.Builder
	_, err := qs.WriteString(`
	SELECT
		id,
		md5
	FROM _file WHERE md5 IN (`)
	if err != nil {
		return "", err
	}
	for i := 0; i < row; i++ {
		if i != 0 {
			qs.WriteString(", ")
		}
		qs.WriteString("?")
	}
	qs.WriteString(");")
	return qs.String(), err
}

func toAny(s []string) []any {
	c := make([]any, len(s))
	for i, v := range s {
		c[i] = v
	}
	return c
}

func (q *Queries) ListFileMd5(ctx context.Context, md5List []string) (m map[string]*File, err error) {
	placeholders := strings.Repeat("?,", len(md5List)-1) + "?"
	query := fmt.Sprintf("SELECT * FROM _file WHERE md5 IN (%s)", placeholders)
	rows, err := q.db.QueryContext(ctx, query, toAny(md5List)...)
	if err != nil {
		return
	}
	defer rows.Close()
	m = make(map[string]*File)
	for rows.Next() {
		var i File
		err = rows.Scan(
			&i.ID,
			&i.Size,
			&i.Md5,
			&i.Sha256,
		)
		if err != nil {
			return
		}
		m[i.Md5] = &i
	}
	return
}
