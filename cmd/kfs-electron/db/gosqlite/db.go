package gosqlite

import (
	"context"
	"database/sql"
	"errors"
	_ "modernc.org/sqlite"
	"os"
	"path/filepath"
)

type DB struct {
	dataSourceName string
	ch             chan *sql.DB
}

func NewDb(dataSourceName string) (*DB, error) {
	db, err := open(dataSourceName)
	if err != nil {
		return nil, err
	}
	err = db.Create()
	if err != nil {
		return nil, err
	}
	return db, nil
}

func open(dataSourceName string) (*DB, error) {
	conn, err := sql.Open("sqlite", dataSourceName)
	if err != nil {
		return nil, err
	}
	db := &DB{
		ch: make(chan *sql.DB, 1),
	}
	db.ch <- conn
	return db, err
}

func (db *DB) GetConn() *sql.DB {
	return <-db.ch
}

func (db *DB) PutConn(conn *sql.DB) {
	db.ch <- conn
}

func (db *DB) Remove() error {
	conn := db.GetConn()
	defer db.PutConn(conn)
	_, err := conn.Exec(`
	DROP TABLE IF EXISTS _file;
	DROP TABLE IF EXISTS _scan_history;
	`)
	return err
}

func (db *DB) Create() error {
	conn := db.GetConn()
	defer db.PutConn(conn)
	_, err := conn.Exec(`
	CREATE TABLE IF NOT EXISTS _file (
	    taskName    INT64          NOT NULL,
	    dirname     VARCHAR(32767) NOT NULL,
		name        VARCHAR(255)   NOT NULL,
		hash        CHAR(64)       NOT NULL,
	    mode        INT64          NOT NULL,
		size        INT64          NOT NULL,
		modifyTime  INT64          NOT NULL,
		PRIMARY KEY (taskName, dirname, name),
		FOREIGN KEY (taskName)     REFERENCES _backup_task(name)
	);
	CREATE TABLE IF NOT EXISTS _backup_task (
		name        VARCHAR(256)   NOT NULL PRIMARY KEY,
		description VARCHAR(256)   NOT NULL,
		srcPath     VARCHAR(32767) NOT NULL,
		driverName  VARCHAR(256)   NOT NULL,
		dstPath     VARCHAR(32767) NOT NULL,
		encoder     VARCHAR(64)    NOT NULL,
	    concurrent  INT8           NOT NULL
	);
	`)
	return err
}

func (db *DB) Close() error {
	select {
	case conn := <-db.ch:
		return conn.Close()
	default:
		return nil
	}
}

func (db *DB) UpsertFile(ctx context.Context, taskName string, path string, hash string, mode os.FileMode, size int64, modifyTime int64) error {
	conn := db.GetConn()
	defer db.PutConn(conn)
	_, err := conn.ExecContext(ctx, `
	INSERT INTO _file VALUES (?, ?, ?, ?, ?, ?) ON CONFLICT(taskName, dirname, name) DO UPDATE SET
		hash=?,
		mode=?,
		size=?,
		modifyTime=?;
	`, taskName, filepath.Dir(path), filepath.Base(path), hash, mode, size, modifyTime, hash, mode, size, modifyTime)
	return err
}

func (db *DB) UpsertBackupTask(ctx context.Context, name string, description string, srcPath string, driverName string, dstPath string, encoder string, concurrent int) error {
	conn := db.GetConn()
	defer db.PutConn(conn)
	_, err := conn.ExecContext(ctx, `
	INSERT INTO _backup_task VALUES (?, ?, ?, ?, ?, ?, ?) ON CONFLICT(name) DO UPDATE SET
		description=?,
		srcPath=?,
		driverName=?,
		dstPath=?,
		encoder=?,
		concurrent=?;
	`, name, description, srcPath, driverName, dstPath, encoder, concurrent, description, srcPath, driverName, dstPath, encoder, concurrent)
	return err
}

func (db *DB) DeleteBackupTask(ctx context.Context, name string) error {
	conn := db.GetConn()
	defer db.PutConn(conn)
	_, err := conn.ExecContext(ctx, `
	DELETE FROM _backup_task WHERE name = ?`, name)
	if err != nil {
		return err
	}
	return err
}

func (db *DB) GetBackupTask(ctx context.Context, name string) (t BackupTask, err error) {
	conn := db.GetConn()
	defer db.PutConn(conn)
	rows, err := conn.QueryContext(ctx, `
	SELECT  description,
			srcPath,
			driverName,
			dstPath,
			encoder,
			concurrent
	FROM _backup_task WHERE name=?;
	`, name)
	if err != nil {
		return
	}
	defer rows.Close()
	if !rows.Next() {
		err = errors.New("no such backup task: " + name)
		return
	}
	t.Name = name
	err = rows.Scan(&t.Description, &t.SrcPath, &t.DriverName, &t.DstPath, &t.Encoder, &t.Concurrent)
	if err != nil {
		return
	}
	return
}

type BackupTask struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	SrcPath     string `json:"srcPath"`
	DriverName  string `json:"driverName"`
	DriverId    uint64 `json:"driverId"`
	DstPath     string `json:"dstPath"`
	Encoder     string `json:"encoder"`
	Concurrent  int    `json:"concurrent"`
}

func (db *DB) ListBackupTask(ctx context.Context) (list []BackupTask, err error) {
	conn := db.GetConn()
	defer db.PutConn(conn)
	rows, err := conn.QueryContext(ctx, `
	SELECT * FROM _backup_task;
	`)
	if err != nil {
		return
	}
	defer rows.Close()
	list = []BackupTask{}
	for rows.Next() {
		var t BackupTask
		err = rows.Scan(&t.Name, &t.Description, &t.SrcPath, &t.DriverName, &t.DstPath, &t.Encoder, &t.Concurrent)
		if err != nil {
			return
		}
		list = append(list, t)
	}
	return
}
