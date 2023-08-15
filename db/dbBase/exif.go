package dbBase

import (
	"context"
	"database/sql"
	"github.com/lazyxu/kfs/dao"
)

func InsertNullExif(ctx context.Context, conn *sql.DB, db DbImpl, hash string) (exist bool, err error) {
	_, err = conn.ExecContext(ctx, `
	INSERT INTO _exif (
		hash  
	) VALUES (?)`, hash)
	if db.IsUniqueConstraintError(err) {
		exist = true
		err = nil
	}
	return
}

func InsertExif(ctx context.Context, conn *sql.DB, db DbImpl, hash string, e dao.Exif) (exist bool, err error) {
	_, err = conn.ExecContext(ctx, `
	INSERT INTO _exif (
		hash,
		version,
		dateTime,
		hostComputer,
		GPSLatitudeRef,
		GPSLatitude,
		GPSLongitudeRef,
		GPSLongitude
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`, hash, e.Version, e.DateTime, e.HostComputer, e.GPSLatitudeRef, e.GPSLatitude, e.GPSLongitudeRef, e.GPSLongitude)
	if db.IsUniqueConstraintError(err) {
		exist = true
		err = nil
	}
	return
}
func ListExpectExif(ctx context.Context, conn *sql.DB) (hashList []string, err error) {
	rows, err := conn.QueryContext(ctx, `
	SELECT hash FROM _file EXCEPT SELECT hash FROM _exif;
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

func ListExpectExifCb(ctx context.Context, conn *sql.DB, cb func(hash string)) (err error) {
	rows, err := conn.QueryContext(ctx, `
	SELECT hash FROM _file EXCEPT SELECT hash FROM _exif;
	`)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var hash string
		err = rows.Scan(&hash)
		if err != nil {
			return
		}
		cb(hash)
	}
	return
}

func ListExif(ctx context.Context, conn *sql.DB) (exifMap map[string]dao.Exif, err error) {
	rows, err := conn.QueryContext(ctx, `
	SELECT 
		hash,
		version,
		dateTime,
		hostComputer,
		GPSLatitudeRef,
		GPSLatitude,
		GPSLongitudeRef,
		GPSLongitude
	FROM _exif WHERE version IS NOT NULL;
	`)
	if err != nil {
		return
	}
	defer rows.Close()
	exifMap = make(map[string]dao.Exif)
	for rows.Next() {
		var hash string
		var exif dao.Exif
		err = rows.Scan(&hash, &exif.Version, &exif.DateTime, &exif.HostComputer,
			&exif.GPSLatitudeRef, &exif.GPSLatitude, &exif.GPSLongitudeRef, &exif.GPSLongitude)
		if err != nil {
			return
		}
		exifMap[hash] = exif
	}
	return
}
