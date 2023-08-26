package dbBase

import (
	"context"
	"database/sql"
	"errors"

	"github.com/lazyxu/kfs/dao"
)

func InsertNullVideoMetadata(ctx context.Context, conn *sql.DB, db DbImpl, hash string) (exist bool, err error) {
	_, err = conn.ExecContext(ctx, `
	INSERT INTO _video_metadata (
		hash  
	) VALUES (?)`, hash)
	if db.IsUniqueConstraintError(err) {
		exist = true
		err = nil
	}
	return
}

func InsertVideoMetadata(ctx context.Context, conn *sql.DB, db DbImpl, hash string, m dao.VideoMetadata) (exist bool, err error) {
	_, err = conn.ExecContext(ctx, `
	INSERT INTO _video_metadata (
		hash,
		Codec,
		Created,
		Modified,
		Duration
	) VALUES (?, ?, ?, ?, ?)`, hash, m.Codec, m.Created, m.Modified, m.Duration)
	if db.IsUniqueConstraintError(err) {
		exist = true
		err = nil
	}
	return
}

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

func ListMetadata(ctx context.Context, conn *sql.DB) (list []dao.Metadata, err error) {
	list = make([]dao.Metadata, 0)
	{
		var rows *sql.Rows
		rows, err = conn.QueryContext(ctx, `
		SELECT 
			_file_type.hash,
			_file_type.Type,
			_file_type.SubType,
			_file_type.Extension,
			Codec,
			Created,
			Modified,
			Duration
		FROM _video_metadata LEFT JOIN _file_type WHERE Codec IS NOT NULL AND _video_metadata.hash=_file_type.hash;
		`)
		if err != nil {
			return
		}
		defer rows.Close()
		for rows.Next() {
			var hash string
			var m dao.VideoMetadata
			var t dao.FileType
			err = rows.Scan(&hash, &t.Type, &t.SubType, &t.Extension,
				&m.Codec, &m.Created, &m.Modified, &m.Duration)
			if err != nil {
				return
			}
			list = append(list, dao.Metadata{
				Hash:          hash,
				FileType:      t,
				VideoMetadata: m,
			})
		}
	}
	{
		var rows *sql.Rows
		rows, err = conn.QueryContext(ctx, `
		SELECT 
			_file_type.hash,
			_file_type.Type,
			_file_type.SubType,
			_file_type.Extension,
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
		FROM _exif LEFT JOIN _file_type WHERE exifVersion IS NOT NULL AND _exif.hash=_file_type.hash;
		`)
		if err != nil {
			return
		}
		defer rows.Close()
		for rows.Next() {
			var hash string
			var e dao.Exif
			var t dao.FileType
			err = rows.Scan(&hash, &t.Type, &t.SubType, &t.Extension,
				&e.ExifVersion, &e.ImageDescription, &e.Orientation,
				&e.DateTime, &e.DateTimeOriginal, &e.DateTimeDigitized,
				&e.OffsetTime, &e.OffsetTimeOriginal, &e.OffsetTimeDigitized,
				&e.SubsecTime, &e.SubsecTimeOriginal, &e.SubsecTimeDigitized,
				&e.HostComputer, &e.Make, &e.Model,
				&e.ExifImageWidth, &e.ExifImageLength,
				&e.GPSLatitudeRef, &e.GPSLatitude, &e.GPSLongitudeRef, &e.GPSLongitude)
			if err != nil {
				return
			}
			list = append(list, dao.Metadata{
				Hash:     hash,
				FileType: t,
				Exif:     e,
			})
		}
	}
	return
}

var ErrNoRecords = errors.New("no such records in db")

func GetMetadata(ctx context.Context, conn *sql.DB, hash string) (metadata dao.Metadata, err error) {
	rows, err := conn.QueryContext(ctx, `
	SELECT 
		_exif.hash,
		_file_type.Type,
		_file_type.SubType,
		_file_type.Extension,
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
	FROM _exif LEFT JOIN _file_type WHERE  _exif.hash=? AND _file_type.hash=? AND _exif.exifVersion IS NOT NULL;
	`, hash, hash)
	if err != nil {
		return
	}
	defer rows.Close()
	if rows.Next() {
		var hash string
		var e dao.Exif
		var t dao.FileType
		err = rows.Scan(&hash, &t.Type, &t.SubType, &t.Extension,
			&e.ExifVersion, &e.ImageDescription, &e.Orientation,
			&e.DateTime, &e.DateTimeOriginal, &e.DateTimeDigitized,
			&e.OffsetTime, &e.OffsetTimeOriginal, &e.OffsetTimeDigitized,
			&e.SubsecTime, &e.SubsecTimeOriginal, &e.SubsecTimeDigitized,
			&e.HostComputer, &e.Make, &e.Model,
			&e.ExifImageWidth, &e.ExifImageLength,
			&e.GPSLatitudeRef, &e.GPSLatitude, &e.GPSLongitudeRef, &e.GPSLongitude)
		if err != nil {
			return
		}
		metadata = dao.Metadata{
			Hash:     hash,
			Exif:     e,
			FileType: t,
		}
	}
	return
}
