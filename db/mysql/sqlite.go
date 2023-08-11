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

	CREATE TABLE IF NOT EXISTS _dir (
		hash       CHAR(64) NOT NULL PRIMARY KEY,
		size       BIGINT  NOT NULL,
		count      BIGINT  NOT NULL,
		totalCount BIGINT  NOT NULL
	);

	CREATE TABLE IF NOT EXISTS _dirItem (
		hash           CHAR(64)     NOT NULL,
		itemHash       CHAR(64)     NOT NULL,
		itemName       VARCHAR(256) NOT NULL,
		itemMode       BIGINT       NOT NULL,
		itemSize       BIGINT       NOT NULL,
		itemCount      BIGINT       NOT NULL,
		itemTotalCount BIGINT       NOT NULL,
		itemCreateTime BIGINT    NOT NULL,
		itemModifyTime BIGINT    NOT NULL,
		itemChangeTime BIGINT    NOT NULL,
		itemAccessTime BIGINT    NOT NULL,
		PRIMARY KEY (Hash, itemName)
	);

	CREATE TABLE IF NOT EXISTS _commit (
		id          BIGINT    NOT NULL PRIMARY KEY AUTO_INCREMENT,
		createTime  BIGINT    NOT NULL,
		Hash        CHAR(64)  NOT NULL,
		lastId      BIGINT    NOT NULL
	);

	CREATE TABLE IF NOT EXISTS _branch (
		name        VARCHAR(256) NOT NULL PRIMARY KEY,
		description VARCHAR(256) NOT NULL DEFAULT "",
		commitId    BIGINT       NOT NULL,
		size        BIGINT       NOT NULL,
		count       BIGINT       NOT NULL,
		FOREIGN KEY (commitId)   REFERENCES _commit(id)
	);

	CREATE TABLE IF NOT EXISTS _branch (
		name        VARCHAR(256) NOT NULL PRIMARY KEY,
		description VARCHAR(256) NOT NULL DEFAULT "",
		commitId    BIGINT       NOT NULL,
		size        BIGINT       NOT NULL,
		count       BIGINT       NOT NULL,
		FOREIGN KEY (commitId)   REFERENCES _commit(id)
	);

	CREATE TABLE IF NOT EXISTS _driver (
		name        VARCHAR(256) NOT NULL PRIMARY KEY,
		description VARCHAR(256) NOT NULL DEFAULT ""
	);

	CREATE TABLE IF NOT EXISTS _driver_file (
		driver_name VARCHAR(256)   NOT NULL,
		filepath    VARCHAR(65536) NOT NULL,
	    version     BIGINT         NOT NULL,
		hash        CHAR(64)       NOT NULL,
		mode        BIGINT         NOT NULL,
		size        BIGINT         NOT NULL,
		createTime  BIGINT         NOT NULL,
		modifyTime  BIGINT         NOT NULL,
		changeTime  BIGINT         NOT NULL,
		accessTime  BIGINT         NOT NULL,
		PRIMARY KEY (driver_name, filepath, version),
		FOREIGN KEY (driver_name)    REFERENCES _driver(name)
	);
	`)
	return err
}

func (db *DB) Close() error {
	//return db.Close()
	return nil
}
