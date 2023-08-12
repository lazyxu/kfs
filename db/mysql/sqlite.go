package mysql

import (
	"database/sql"
	"sync"

	"github.com/lazyxu/kfs/dao"

	_ "github.com/go-sql-driver/mysql"
)

type DB struct {
	db          *sql.DB
	branchCache sync.Map
}

func New(dataSourceName string) (dao.Database, error) {
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
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(10)
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	d := &DB{
		db: db,
	}
	return d, err
}

func (db *DB) Remove() error {
	_, err := db.db.Exec(`
	DROP TABLE IF EXISTS _file, _dir, _dirItem, _commit, _branch, _driver, _files;
	`)
	return err
}

func (db *DB) Create() error {
	_, err := db.db.Exec(`
	CREATE TABLE IF NOT EXISTS _file (
		hash CHAR(64) NOT NULL PRIMARY KEY,
		size BIGINT   NOT NULL
	);

	CREATE TABLE IF NOT EXISTS _driver (
		name        VARCHAR(256) NOT NULL PRIMARY KEY,
		description VARCHAR(256) NOT NULL DEFAULT ""
	);

	CREATE TABLE IF NOT EXISTS _driver_file (
		driver_name VARCHAR(256)   NOT NULL,
		dirpath     VARCHAR(32767) NOT NULL,
		name        VARCHAR(255)   NOT NULL,
	    version     BIGINT         NOT NULL,
		hash        CHAR(64)       NOT NULL,
		mode        BIGINT         NOT NULL,
		size        BIGINT         NOT NULL,
		createTime  BIGINT         NOT NULL,
		modifyTime  BIGINT         NOT NULL,
		changeTime  BIGINT         NOT NULL,
		accessTime  BIGINT         NOT NULL,
		PRIMARY KEY (driver_name, dirpath, name, version),
		FOREIGN KEY (driver_name)  REFERENCES _driver(name)
	);
	`) // https://blog.csdn.net/jimmyleeee/article/details/124682486
	return err
}

func (db *DB) Close() error {
	//return db.Close()
	return nil
}
