package dbBase

import (
	"context"
	"database/sql"
	"github.com/lazyxu/kfs/dao"
)

func InsertDevice(ctx context.Context, conn *sql.DB, id string, name string, os string, userAgent string, hostname string) error {
	_, err := conn.ExecContext(ctx, `
	INSERT INTO _device (
		id,
		name,
		os,
		userAgent,
		hostname
	) VALUES (?, ?, ?, ?, ?)
	ON CONFLICT(id) DO UPDATE SET
		name=?,
		os=?,
		userAgent=?,
		hostname=?`, id, name, os, userAgent, hostname, name, os, userAgent, hostname)
	if err != nil {
		return err
	}
	return nil
}

func DeleteDevice(ctx context.Context, conn *sql.DB, deviceId string) error {
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
		err = rows.Scan(&item.Id, &item.Name, &item.OS, &item.UserAgent, &item.Hostname)
		if err != nil {
			return
		}
		list = append(list, item)
	}
	return
}
