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

func FuncNew(dataSourceName string) func() (dao.DB, error) {
	return func() (dao.DB, error) {
		db, err := Open(dataSourceName)
		if err != nil {
			return nil, err
		}
		err = db.Create()
		if err != nil {
			return nil, err
		}
		return db, nil
	}
}

func Open(dataSourceName string) (*DB, error) {
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

func (db *DB) getConn() *sql.DB {
	return db.db
}

func (db *DB) putConn(conn *sql.DB) {
}

func (db *DB) Remove() error {
	conn := db.getConn()
	defer db.putConn(conn)
	_, err := conn.Exec(`
	DROP TABLE IF EXISTS _file, _dir, _dirItem, _commit, _branch;
	`)
	return err
}

func (db *DB) Create() error {
	conn := db.getConn()
	defer db.putConn(conn)
	_, err := conn.Exec(`
	CREATE TABLE IF NOT EXISTS _file (
		hash CHAR(64) NOT NULL PRIMARY KEY,
		size BIGINT  NOT NULL
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
		itemCreateTime TIMESTAMP    NOT NULL,
		itemModifyTime TIMESTAMP    NOT NULL,
		itemChangeTime TIMESTAMP    NOT NULL,
		itemAccessTime TIMESTAMP    NOT NULL,
		PRIMARY KEY(Hash, itemName)
	);

	CREATE TABLE IF NOT EXISTS _commit (
		id          BIGINT    NOT NULL PRIMARY KEY AUTO_INCREMENT,
		createTime  TIMESTAMP NOT NULL,
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
	`)
	return err
}

func (db *DB) Close() error {
	//return db.Close()
	return nil
}
