package mysql

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type DB struct {
	db *sql.DB
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

func (db *DB) Create() error {
	conn := db.getConn()
	defer db.putConn(conn)
	_, err := conn.Exec(`
	DROP TABLE IF EXISTS _file;
	CREATE TABLE _file (
		hash CHAR(64) NOT NULL PRIMARY KEY,
		size INTEGER  NOT NULL
	);

	DROP TABLE IF EXISTS _dir;
	CREATE TABLE _dir (
		hash       CHAR(64) NOT NULL PRIMARY KEY,
		size       INTEGER  NOT NULL,
		count      INTEGER  NOT NULL,
		totalCount INTEGER  NOT NULL
	);

	DROP TABLE IF EXISTS _dirItem;
	CREATE TABLE _dirItem (
		hash           CHAR(64)     NOT NULL,
		itemHash       CHAR(64)     NOT NULL,
		itemName       VARCHAR(256) NOT NULL,
		itemMode       INTEGER      NOT NULL,
		itemSize       INTEGER      NOT NULL,
		itemCount      INTEGER      NOT NULL,
		itemTotalCount INTEGER      NOT NULL,
		itemCreateTime TIMESTAMP    NOT NULL,
		itemModifyTime TIMESTAMP    NOT NULL,
		itemChangeTime TIMESTAMP    NOT NULL,
		itemAccessTime TIMESTAMP    NOT NULL,
		PRIMARY KEY(Hash, itemName)
	);

	DROP TABLE IF EXISTS _commit;
	CREATE TABLE _commit (
		id          INTEGER   NOT NULL PRIMARY KEY AUTO_INCREMENT,
		createTime  TIMESTAMP NOT NULL,
		Hash        CHAR(64)  NOT NULL,
		lastId      INTEGER   NOT NULL
	);

	DROP TABLE IF EXISTS _branch;
	CREATE TABLE _branch (
		name        VARCHAR(256) NOT NULL PRIMARY KEY,
		description VARCHAR(256) NOT NULL DEFAULT "",
		commitId    INTEGER      NOT NULL,
		size        INTEGER      NOT NULL,
		count       INTEGER      NOT NULL
	);
	`)
	return err
}

func (db *DB) Close() error {
	//return db.Close()
	return nil
}
