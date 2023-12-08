package dbBase

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/lazyxu/kfs/dao"
)

func InsertDriver(ctx context.Context, conn *sql.DB, db DbImpl, driverName string, description string) (exist bool, err error) {
	_, err = conn.ExecContext(ctx, `
	INSERT INTO _driver (
		name,
		description,
	    Type
	) VALUES (?, ?, ?)`, driverName, description, "")
	if db.IsUniqueConstraintError(err) {
		exist = true
		err = nil
	}
	return
}

const DRIVER_TYPE_BAIDU_PHOTO = "baiduPhoto"
const DRIVER_TYPE_LOCAL_FILE = "localFile"

func InsertDriverBaiduPhoto(ctx context.Context, conn *sql.DB, db DbImpl, driverName string, description string, accessToken string, refreshToken string) (exist bool, err error) {
	res, err := conn.ExecContext(ctx, `
	INSERT INTO _driver (
		name,
		description,
	    Type
	) VALUES (?, ?, ?)`, driverName, description, DRIVER_TYPE_BAIDU_PHOTO)
	if db.IsUniqueConstraintError(err) {
		exist = true
		err = nil
	}
	id, err := res.LastInsertId()
	if err != nil {
		return
	}
	_, err = conn.ExecContext(ctx, `
	INSERT INTO _driver_sync (
		id
	) VALUES (?)`, id)
	if err != nil {
		return
	}
	_, err = conn.ExecContext(ctx, `
	INSERT INTO _driver_baidu_photo (
		id,
		accessToken,
	    refreshToken
	) VALUES (?, ?, ?)`, id, accessToken, refreshToken)
	return
}

func InsertDriverLocalFile(ctx context.Context, conn *sql.DB, db DbImpl, driverName string, description string, deviceId uint64, srcPath string, ignores string, encoder string) (exist bool, err error) {
	res, err := conn.ExecContext(ctx, `
	INSERT INTO _driver (
		name,
		description,
	    Type
	) VALUES (?, ?, ?)`, driverName, description, DRIVER_TYPE_LOCAL_FILE)
	if db.IsUniqueConstraintError(err) {
		exist = true
		err = nil
	}
	id, err := res.LastInsertId()
	if err != nil {
		return
	}
	_, err = conn.ExecContext(ctx, `
	INSERT INTO _driver_sync (
		id
	) VALUES (?)`, id)
	if err != nil {
		return
	}
	_, err = conn.ExecContext(ctx, `
	INSERT INTO _driver_local_file (
		id,
		deviceId,
	    srcPath,
	    ignores,
	    encoder
	) VALUES (?, ?, ?, ?, ?)`, id, deviceId, srcPath, ignores, encoder)
	return
}

func UpdateDriverSync(ctx context.Context, conn *sql.DB, driverId uint64, sync bool, h int64, m int64) error {
	_, err := conn.ExecContext(ctx, `
	UPDATE _driver_sync
	SET sync = ?, h = ?, m = ?
	WHERE id = ?;`, sync, h, m, driverId)
	if err != nil {
		return err
	}
	return nil
}

func UpdateDriverLocalFile(ctx context.Context, conn *sql.DB, driverId uint64, srcPath, ignores, encoder string) error {
	_, err := conn.ExecContext(ctx, `
	UPDATE _driver_local_file
	SET srcPath = ?, ignores = ?, encoder = ?
	WHERE id = ?;`, srcPath, ignores, encoder, driverId)
	if err != nil {
		return err
	}
	return nil
}

func ResetDriver(ctx context.Context, conn *sql.DB, driverId uint64) error {
	_, err := conn.ExecContext(ctx, `
	DELETE FROM _driver_file WHERE driverId = ?`, driverId)
	if err != nil {
		return err
	}
	return err
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

func GetDriver(ctx context.Context, txOrDb TxOrDb, driverId uint64) (driver dao.Driver, err error) {
	rows, err := txOrDb.QueryContext(ctx, `
	SELECT * FROM _driver WHERE id = ?;
	`, driverId)
	if err != nil {
		return
	}
	defer rows.Close()
	if rows.Next() {
		err = rows.Scan(&driver.Id, &driver.Name, &driver.Description, &driver.Typ)
		if err != nil {
			return
		}
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
	SELECT sync, h, m FROM _driver_sync WHERE id = ?;
	`, driverId)
	if err != nil {
		return
	}
	defer rows.Close()
	if rows.Next() {
		err = rows.Scan(&driver.Sync, &driver.H, &driver.M)
		if err != nil {
			return
		}
	}
	return
}

func ListCloudDriverSync(ctx context.Context, txOrDb TxOrDb) (drivers []dao.Driver, err error) {
	rows, err := txOrDb.QueryContext(ctx, `
	SELECT _driver.id, h, m FROM _driver_sync LEFT JOIN _driver WHERE _driver_sync.sync = 1 AND _driver_sync.id = _driver.id AND _driver.type = ?;
	`, DRIVER_TYPE_BAIDU_PHOTO)
	if err != nil {
		return
	}
	defer rows.Close()
	drivers = []dao.Driver{}
	for rows.Next() {
		var driver dao.Driver
		err = rows.Scan(&driver.Id, &driver.H, &driver.M)
		if err != nil {
			return
		}
		drivers = append(drivers, driver)
	}
	return
}

func ListLocalFileDriver(ctx context.Context, txOrDb TxOrDb, deviceId uint64) (drivers []dao.Driver, err error) {
	rows, err := txOrDb.QueryContext(ctx, `
	SELECT _driver.id, h, m, srcPath, ignores, encoder
	FROM _driver_local_file LEFT JOIN _driver_sync LEFT JOIN _driver
	WHERE _driver_local_file.deviceId = ? AND _driver_local_file.id = _driver_sync.id AND _driver_sync.sync = 1 AND _driver_sync.id = _driver.id AND _driver.type = ?;
	`, deviceId, DRIVER_TYPE_LOCAL_FILE)
	if err != nil {
		return
	}
	defer rows.Close()
	drivers = []dao.Driver{}
	for rows.Next() {
		var driver dao.Driver
		err = rows.Scan(&driver.Id, &driver.H, &driver.M, &driver.SrcPath, &driver.Ignores, &driver.Encoder)
		if err != nil {
			return
		}
		drivers = append(drivers, driver)
	}
	return
}

func GetDriverLocalFile(ctx context.Context, txOrDb TxOrDb, driverId uint64) (driver dao.Driver, err error) {
	rows, err := txOrDb.QueryContext(ctx, `
	SELECT deviceId, srcPath, ignores, encoder FROM _driver_local_file WHERE id = ?;
	`, driverId)
	if err != nil {
		return
	}
	defer rows.Close()
	if rows.Next() {
		err = rows.Scan(&driver.DeviceId, &driver.SrcPath, &driver.Ignores, &driver.Encoder)
		if err != nil {
			return
		}
	}
	return
}

func GetDriverFileSize(ctx context.Context, txOrDb TxOrDb, driverId uint64) (n uint64, err error) {
	rows, err := txOrDb.QueryContext(ctx, `
	SELECT IFNULL(SUM(size), 0) FROM _driver_file WHERE driverId = ? AND mode < 2147483648;
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

func GetDriverDirCalculatedInfo(ctx context.Context, txOrDb TxOrDb, driverId uint64, filePath []string) (info dao.DirCalculatedInfo, err error) {
	var like string
	if len(filePath) != 0 {
		data, err := json.Marshal(filePath)
		if err != nil {
			panic(err)
		}
		like = string(data)
		like = like[:len(like)-1] + "%"
	} else {
		like = "%"
	}
	var rows *sql.Rows
	{
		rows, err = txOrDb.QueryContext(ctx, `
	SELECT COUNT(1), IFNULL(SUM(size), 0) FROM _driver_file WHERE driverId = ? AND mode < 2147483648 AND dirPath LIKE ?;
	`, driverId, like)
		if err != nil {
			return
		}
		if rows.Next() {
			err = rows.Scan(&info.FileCount, &info.FileSize)
			if err != nil {
				return
			}
		}
		rows.Close()
	}
	{
		rows, err = txOrDb.QueryContext(ctx, `
	SELECT COUNT(1), IFNULL(SUM(size), 0) FROM (SELECT distinct hash, size FROM _driver_file WHERE driverId = ? AND mode < 2147483648 AND dirPath LIKE ?)
	`, driverId, like)
		if err != nil {
			return
		}
		if rows.Next() {
			err = rows.Scan(&info.DistinctFileCount, &info.DistinctFileSize)
			if err != nil {
				return
			}
		}
		rows.Close()
	}
	{
		rows, err = txOrDb.QueryContext(ctx, `
	SELECT COUNT(1) FROM _driver_file WHERE driverId = ? AND mode >= 2147483648 AND dirPath LIKE ?;
	`, driverId, like)
		if err != nil {
			return
		}
		if rows.Next() {
			err = rows.Scan(&info.DirCount)
			if err != nil {
				return
			}
		}
		rows.Close()
	}
	return
}
