package dbBase

import (
	"context"
	"database/sql"
	"github.com/lazyxu/kfs/dao"
)

func InsertDriver(ctx context.Context, conn *sql.DB, db DbImpl, driverName string, description string, typ string, accessToken string, refreshToken string) (exist bool, err error) {
	_, err = conn.ExecContext(ctx, `
	INSERT INTO _driver (
		name,
		description,
	    Type,
		accessToken,
		refreshToken
	) VALUES (?, ?, ?, ?, ?)`, driverName, description, typ, accessToken, refreshToken)
	if db.IsUniqueConstraintError(err) {
		exist = true
		err = nil
	}
	return
}

func UpdateDriverSync(ctx context.Context, conn *sql.DB, driverName string, sync bool, h int64, m int64, s int64) error {
	_, err := conn.ExecContext(ctx, `
	UPDATE _driver
	SET sync = ?, h = ?, m = ?, s = ?
	WHERE name = ?;`, sync, h, m, s, driverName)
	if err != nil {
		return err
	}
	return nil
}

func DeleteDriver(ctx context.Context, conn *sql.DB, driverName string) error {
	_, err := conn.ExecContext(ctx, `
	DELETE FROM _driver WHERE name = ?`, driverName)
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
		err = rows.Scan(&driver.Name, &driver.Description, &driver.Typ, &driver.Sync, &driver.H, &driver.M, &driver.S, &driver.AccessToken, &driver.RefreshToken)
		if err != nil {
			return
		}
		drivers = append(drivers, driver)
	}
	return
}

func GetDriverToken(ctx context.Context, txOrDb TxOrDb, driverName string) (driver dao.Driver, err error) {
	rows, err := txOrDb.QueryContext(ctx, `
	SELECT accessToken, refreshToken FROM _driver WHERE name = ?;
	`, driverName)
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

func GetDriverSync(ctx context.Context, txOrDb TxOrDb, driverName string) (driver dao.Driver, err error) {
	rows, err := txOrDb.QueryContext(ctx, `
	SELECT sync, h, m, s FROM _driver WHERE name = ?;
	`, driverName)
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

func GetDriverFileSize(ctx context.Context, txOrDb TxOrDb, driverName string) (n uint64, err error) {
	rows, err := txOrDb.QueryContext(ctx, `
	SELECT IFNULL(SUM(size), 0) FROM _driver_file WHERE driverName = ? AND mode < 2147483648;;
	`, driverName)
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

func GetDriverFileCount(ctx context.Context, txOrDb TxOrDb, driverName string) (n uint64, err error) {
	rows, err := txOrDb.QueryContext(ctx, `
	SELECT COUNT(1) FROM _driver_file WHERE driverName = ? AND mode < 2147483648;
	`, driverName)
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

func GetDriverDirCount(ctx context.Context, txOrDb TxOrDb, driverName string) (n uint64, err error) {
	rows, err := txOrDb.QueryContext(ctx, `
	SELECT COUNT(1) FROM _driver_file WHERE driverName = ? AND mode >= 2147483648;
	`, driverName)
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
