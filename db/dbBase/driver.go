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
		err = rows.Scan(&driver.Name, &driver.Description, &driver.Typ, &driver.AccessToken, &driver.RefreshToken)
		if err != nil {
			return
		}
		drivers = append(drivers, driver)
	}
	return
}
