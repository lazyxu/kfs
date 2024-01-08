package dbBase

import (
	"context"
	"github.com/lazyxu/kfs/dao"
)

func ListDCIMDriverMetadata(ctx context.Context, txOrDb TxOrDb, driver *dao.DCIMDriver) error {
	rows, err := txOrDb.QueryContext(ctx, `
		SELECT 
			_file_type.hash,
			_file_type.Type,
			_file_type.SubType,
			_file_type.Extension,
			time,
			year,
			month,
			day
		FROM (
		    SELECT
		        DISTINCT hash
		    FROM _driver_file WHERE driverId = ? AND mode < 2147483648
		) AS table1 LEFT JOIN _dcim_metadata_time LEFT JOIN _file_type WHERE table1.hash=_dcim_metadata_time.hash AND _dcim_metadata_time.hash=_file_type.hash ORDER BY time DESC;
	`, driver.Id)
	if err != nil {
		return err
	}
	defer rows.Close()
	driver.MetadataList = make([]dao.Metadata, 0)
	for rows.Next() {
		m := dao.Metadata{FileType: &dao.FileType{}}
		err = rows.Scan(&m.Hash, &m.FileType.Type, &m.FileType.SubType, &m.FileType.Extension,
			&m.Time, &m.Year, &m.Month, &m.Day)
		if err != nil {
			return err
		}
		driver.MetadataList = append(driver.MetadataList, m)
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
