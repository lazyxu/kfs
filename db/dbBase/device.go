package dbBase

import (
	"context"
	"database/sql"
	"github.com/lazyxu/kfs/dao"
)

func InsertDevice(ctx context.Context, conn *sql.DB, name string, os string) (int64, error) {
	res, err := conn.ExecContext(ctx, `
	INSERT INTO _device (
		name,
		os
	) VALUES (?, ?)`, name, os)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func DeleteDevice(ctx context.Context, conn *sql.DB, deviceId uint64) error {
	_, err := conn.ExecContext(ctx, `
	DELETE FROM _device WHERE id = ?`, deviceId)
	if err != nil {
		return err
	}
	return err
}

func ListDevice(ctx context.Context, conn *sql.DB) (list []dao.Device, err error) {
	rows, err := conn.QueryContext(ctx, `
	SELECT * FROM _device;
	`)
	if err != nil {
		return
	}
	defer rows.Close()
	list = []dao.Device{}
	for rows.Next() {
		var item dao.Device
		err = rows.Scan(&item.Id, &item.Name, &item.OS)
		if err != nil {
			return
		}
		list = append(list, item)
	}
	return
}
