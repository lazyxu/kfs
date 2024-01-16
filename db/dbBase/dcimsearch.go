package dbBase

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/lazyxu/kfs/dao"
	"strings"
)

func SearchDCIM(ctx context.Context, conn *sql.DB, typeList []string, suffixList []string) (list []dao.Metadata, err error) {
	list = make([]dao.Metadata, 0)
	var whereList []string
	var typeWhere string
	for i, typ := range typeList {
		l := strings.Split(typ, "/")
		if len(l) != 2 {
			err = fmt.Errorf("invalid typeList: %+v\n", typeList)
			return
		}
		t := l[0]
		subT := l[1]
		where := fmt.Sprintf("_file_type.Type='%s' AND _file_type.SubType='%s'", t, subT)
		if i != 0 {
			typeWhere += " OR "
		}
		typeWhere += where
	}
	if typeWhere != "" {
		whereList = append(whereList, typeWhere)
	}
	var suffixWhere string
	for i, suffix := range suffixList {
		where := fmt.Sprintf("_file_type.Extension='%s'", suffix)
		if i != 0 {
			suffixWhere += " OR "
		}
		suffixWhere += where
	}
	if suffixWhere != "" {
		whereList = append(whereList, suffixWhere)
	}
	var where string
	for _, w := range whereList {
		where += " AND (" + w + ")"
	}
	var query = `
		SELECT 
			_file_type.hash,
			_file_type.Type,
			_file_type.SubType,
			_file_type.Extension,
			_height_width.height,
			_height_width.width,
			time,
			year,
			month,
			day,
			IFNULL(_video_metadata.duration, -1)
		FROM _dcim_metadata_time INNER JOIN _height_width INNER JOIN _file_type LEFT JOIN _video_metadata
		ON _dcim_metadata_time.hash=_video_metadata.hash
		WHERE _dcim_metadata_time.hash=_height_width.hash AND _dcim_metadata_time.hash=_file_type.hash` + where + `
		ORDER BY time DESC;
		`
	var rows *sql.Rows
	rows, err = conn.QueryContext(ctx, query)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		md := dao.Metadata{FileType: &dao.FileType{}, HeightWidth: &dao.HeightWidth{}}
		err = rows.Scan(&md.Hash, &md.FileType.Type, &md.FileType.SubType, &md.FileType.Extension,
			&md.HeightWidth.Height, &md.HeightWidth.Width,
			&md.Time, &md.Year, &md.Month, &md.Day, &md.Duration)
		if err != nil {
			return
		}
		list = append(list, md)
	}
	return
}
