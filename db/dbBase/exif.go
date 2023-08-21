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
		ExifVersion,
		ImageDescription,
		Orientation,
		DateTime,
		DateTimeOriginal,
		DateTimeDigitized,
		OffsetTime,
		OffsetTimeOriginal,
		OffsetTimeDigitized,
		SubsecTime,
		SubsecTimeOriginal,
		SubsecTimeDigitized,
		HostComputer,
		Make,
		Model,
		ExifImageWidth,
		ExifImageLength,
		GPSLatitudeRef,
		GPSLatitude,
		GPSLongitudeRef,
		GPSLongitude
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`, hash, e.ExifVersion, e.ImageDescription, e.Orientation,
		e.DateTime, e.DateTimeOriginal, e.DateTimeDigitized,
		e.OffsetTime, e.OffsetTimeOriginal, e.OffsetTimeDigitized,
		e.SubsecTime, e.SubsecTimeOriginal, e.SubsecTimeDigitized,
		e.HostComputer, e.Make, e.Model,
		e.ExifImageWidth, e.ExifImageLength,
		e.GPSLatitudeRef, e.GPSLatitude, e.GPSLongitudeRef, e.GPSLongitude)
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
		ExifVersion,
		ImageDescription,
		Orientation,
		DateTime,
		DateTimeOriginal,
		DateTimeDigitized,
		OffsetTime,
		OffsetTimeOriginal,
		OffsetTimeDigitized,
		SubsecTime,
		SubsecTimeOriginal,
		SubsecTimeDigitized,
		HostComputer,
		Make,
		Model,
		ExifImageWidth,
		ExifImageLength,
		GPSLatitudeRef,
		GPSLatitude,
		GPSLongitudeRef,
		GPSLongitude
	FROM _exif WHERE exifVersion IS NOT NULL;
	`)
	if err != nil {
		return
	}
	defer rows.Close()
	exifMap = make(map[string]dao.Exif)
	for rows.Next() {
		var hash string
		var e dao.Exif
		err = rows.Scan(&hash, &e.ExifVersion, &e.ImageDescription, &e.Orientation,
			&e.DateTime, &e.DateTimeOriginal, &e.DateTimeDigitized,
			&e.OffsetTime, &e.OffsetTimeOriginal, &e.OffsetTimeDigitized,
			&e.SubsecTime, &e.SubsecTimeOriginal, &e.SubsecTimeDigitized,
			&e.HostComputer, &e.Make, &e.Model,
			&e.ExifImageWidth, &e.ExifImageLength,
			&e.GPSLatitudeRef, &e.GPSLatitude, &e.GPSLongitudeRef, &e.GPSLongitude)
		if err != nil {
			return
		}
		exifMap[hash] = e
	}
	return
}
