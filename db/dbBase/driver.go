package dbBase

import (
	"context"
	"database/sql"
	"github.com/lazyxu/kfs/dao"
)

func InsertDriver(ctx context.Context, conn *sql.DB, db DbImpl, driverName string, description string) (exist bool, err error) {
	_, err = conn.ExecContext(ctx, `
	INSERT INTO _driver (
		name,
		description
	) VALUES (?, ?)`, driverName, description)
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

func ListDriver(ctx context.Context, txOrDb TxOrDb) (branches []dao.IDriver, err error) {
	rows, err := txOrDb.QueryContext(ctx, `
	SELECT * FROM _driver;
	`)
	if err != nil {
		return
	}
	defer rows.Close()
	branches = []dao.IDriver{}
	for rows.Next() {
		var branch dao.Branch
		err = rows.Scan(&branch.Name, &branch.Description)
		if err != nil {
			return
		}
		branches = append(branches, branch)
	}
	return
}
