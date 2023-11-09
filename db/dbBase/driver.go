package dbBase

import (
	"context"
	"database/sql"
	"github.com/lazyxu/kfs/dao"
)

func InsertDriver(ctx context.Context, conn *sql.DB, db DbImpl, driverName string, description string, typ string) (exist bool, err error) {
	_, err = conn.ExecContext(ctx, `
	INSERT INTO _driver (
		name,
		description,
	    Type
	) VALUES (?, ?, ?)`, driverName, description, typ)
	if db.IsUniqueConstraintError(err) {
		exist = true
		err = nil
	}
	return
}

func InsertDriverBaiduPhoto(ctx context.Context, conn *sql.DB, db DbImpl, driverName string, description string, typ string, accessToken string, refreshToken string) (exist bool, err error) {
	res, err := conn.ExecContext(ctx, `
	INSERT INTO _driver (
		name,
		description,
	    Type
	) VALUES (?, ?, ?)`, driverName, description, typ)
	if db.IsUniqueConstraintError(err) {
		exist = true
		err = nil
	}
	id, err := res.LastInsertId()
	if err != nil {
		return
	}
	_, err = conn.ExecContext(ctx, `
	INSERT INTO _driver_baidu_photo (
		id,
		accessToken,
	    refreshToken
	) VALUES (?, ?, ?)`, id, accessToken, refreshToken)
	if db.IsUniqueConstraintError(err) {
		exist = true
		err = nil
	}
	return
}

func UpdateDriverSync(ctx context.Context, conn *sql.DB, driverId uint64, sync bool, h int64, m int64, s int64) error {
	_, err := conn.ExecContext(ctx, `
	UPDATE _driver
	SET sync = ?, h = ?, m = ?, s = ?
	WHERE id = ?;`, sync, h, m, s, driverId)
	if err != nil {
		return err
	}
	return nil
}

func DeleteDriver(ctx context.Context, conn *sql.DB, driverId uint64) error {
	_, err := conn.ExecContext(ctx, `
	DELETE FROM _driver WHERE id = ?`, driverId)
	if err != nil {
		return err
	}
	return err
}

func ListDriver(ctx context.Context, txOrDb TxOrDb) (drivers []dao.Driver, err error) {
	rows, err := txOrDb.QueryContext(ctx, `
	SELECT * FROM _driver;
	`)
	if err != nil {
		return
	}
	defer rows.Close()
	drivers = []dao.Driver{}
	for rows.Next() {
		var driver dao.Driver
		err = rows.Scan(&driver.Id, &driver.Name, &driver.Description, &driver.Typ)
		if err != nil {
			return
		}
		drivers = append(drivers, driver)
	}
	return
}

func GetDriverToken(ctx context.Context, txOrDb TxOrDb, driverId uint64) (driver dao.Driver, err error) {
	rows, err := txOrDb.QueryContext(ctx, `
	SELECT accessToken, refreshToken FROM _driver_baidu_photo WHERE id = ?;
	`, driverId)
	if err != nil {
		return
	}
	defer rows.Close()
	if rows.Next() {
		err = rows.Scan(&driver.AccessToken, &driver.RefreshToken)
		if err != nil {
			return
		}
	}
	return
}

func GetDriverSync(ctx context.Context, txOrDb TxOrDb, driverId uint64) (driver dao.Driver, err error) {
	rows, err := txOrDb.QueryContext(ctx, `
	SELECT sync, h, m, s FROM _driver_baidu_photo WHERE id = ?;
	`, driverId)
	if err != nil {
		return
	}
	defer rows.Close()
	if rows.Next() {
		err = rows.Scan(&driver.Sync, &driver.H, &driver.M, &driver.S)
		if err != nil {
			return
		}
	}
	return
}

func GetDriverFileSize(ctx context.Context, txOrDb TxOrDb, driverId uint64) (n uint64, err error) {
	rows, err := txOrDb.QueryContext(ctx, `
	SELECT IFNULL(SUM(size), 0) FROM _driver_file WHERE driverId = ? AND mode < 2147483648;;
	`, driverId)
	if err != nil {
		return
	}
	defer rows.Close()
	if rows.Next() {
		err = rows.Scan(&n)
		if err != nil {
			return
		}
	}
	return
}

func GetDriverFileCount(ctx context.Context, txOrDb TxOrDb, driverId uint64) (n uint64, err error) {
	rows, err := txOrDb.QueryContext(ctx, `
	SELECT COUNT(1) FROM _driver_file WHERE driverId = ? AND mode < 2147483648;
	`, driverId)
	if err != nil {
		return
	}
	defer rows.Close()
	if rows.Next() {
		err = rows.Scan(&n)
		if err != nil {
			return
		}
	}
	return
}

func GetDriverDirCount(ctx context.Context, txOrDb TxOrDb, driverId uint64) (n uint64, err error) {
	rows, err := txOrDb.QueryContext(ctx, `
	SELECT COUNT(1) FROM _driver_file WHERE driverId = ? AND mode >= 2147483648;
	`, driverId)
	if err != nil {
		return
	}
	defer rows.Close()
	if rows.Next() {
		err = rows.Scan(&n)
		if err != nil {
			return
		}
	}
	return
}
