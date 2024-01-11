package dbBase

import (
	"context"
	"database/sql"
	"github.com/lazyxu/kfs/dao"
)

func ListMetadataTime(ctx context.Context, conn *sql.DB) (list []dao.Metadata, err error) {
	list = make([]dao.Metadata, 0)
	var rows *sql.Rows
	rows, err = conn.QueryContext(ctx, `
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
		WHERE _dcim_metadata_time.hash=_height_width.hash AND _dcim_metadata_time.hash=_file_type.hash
		ORDER BY time DESC;
		`)
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

func ListDCIMDriverMetadata(ctx context.Context, txOrDb TxOrDb, driver *dao.DCIMDriver) error {
	rows, err := txOrDb.QueryContext(ctx, `
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
		FROM (
		    SELECT
		        DISTINCT hash
		    FROM _driver_file WHERE driverId = ? AND mode < 2147483648
		) AS table1 INNER JOIN _dcim_metadata_time INNER JOIN _height_width INNER JOIN _file_type LEFT JOIN _video_metadata
		ON _dcim_metadata_time.hash=_video_metadata.hash
		WHERE table1.hash=_dcim_metadata_time.hash AND _dcim_metadata_time.hash=_height_width.hash AND _dcim_metadata_time.hash=_file_type.hash
		ORDER BY time DESC;
	`, driver.Id)
	if err != nil {
		return err
	}
	defer rows.Close()
	driver.MetadataList = make([]dao.Metadata, 0)
	for rows.Next() {
		md := dao.Metadata{FileType: &dao.FileType{}, HeightWidth: &dao.HeightWidth{}}
		err = rows.Scan(&md.Hash, &md.FileType.Type, &md.FileType.SubType, &md.FileType.Extension,
			&md.HeightWidth.Height, &md.HeightWidth.Width,
			&md.Time, &md.Year, &md.Month, &md.Day, &md.Duration)
		if err != nil {
			return err
		}
		driver.MetadataList = append(driver.MetadataList, md)
	}
	return nil
}

func ListDCIMDriver(ctx context.Context, txOrDb TxOrDb) (drivers []dao.DCIMDriver, err error) {
	rows, err := txOrDb.QueryContext(ctx, `
	SELECT * FROM _driver;
	`)
	if err != nil {
		return
	}
	defer rows.Close()
	drivers = []dao.DCIMDriver{}
	for rows.Next() {
		var driver dao.DCIMDriver
		err = rows.Scan(&driver.Id, &driver.Name, &driver.Description, &driver.Typ)
		if err != nil {
			return
		}
		drivers = append(drivers, driver)
	}
	for i := range drivers {
		err = ListDCIMDriverMetadata(ctx, txOrDb, &drivers[i])
		if err != nil {
			return nil, err
		}
	}
	return
}

func ListDCIMMediaType(ctx context.Context, conn *sql.DB) (m map[string][]dao.Metadata, err error) {
	m = make(map[string][]dao.Metadata)
	{
		list := make([]dao.Metadata, 0)
		var rows *sql.Rows
		rows, err = conn.QueryContext(ctx, `
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
		WHERE _dcim_metadata_time.hash=_height_width.hash AND _dcim_metadata_time.hash=_file_type.hash AND _file_type.Type="video"
		ORDER BY time DESC;
		`)
		if err != nil {
			return
		}
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
		rows.Close()
		m["video"] = list
	}
	{
		list := make([]dao.Metadata, 0)
		var rows *sql.Rows
		rows, err = conn.QueryContext(ctx, `
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
		FROM _exif INNER JOIN _dcim_metadata_time INNER JOIN _height_width INNER JOIN _file_type LEFT JOIN _video_metadata
		ON _dcim_metadata_time.hash=_video_metadata.hash
		WHERE ( _exif.Orientation=2 OR _exif.Orientation=4 OR _exif.Orientation=5 OR _exif.Orientation=7)
		  AND _exif.hash=_dcim_metadata_time.hash AND _dcim_metadata_time.hash=_height_width.hash AND _dcim_metadata_time.hash=_file_type.hash
		ORDER BY time DESC;
		`)
		if err != nil {
			return
		}
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
		rows.Close()
		m["selfie"] = list
	}
	return
}
