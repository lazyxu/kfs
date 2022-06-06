package noncgo

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

type DB struct {
	ch chan *sql.DB
}

func Open(dataSourceName string) (*DB, error) {
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

func (db *DB) getConn() *sql.DB {
	return <-db.ch
}

func (db *DB) putConn(conn *sql.DB) {
	db.ch <- conn
}

func (db *DB) Create() error {
	conn := db.getConn()
	defer db.putConn(conn)
	_, err := conn.Exec(`
	DROP TABLE IF EXISTS file;
	CREATE TABLE file (
		hash CHAR(64) NOT NULL PRIMARY KEY,
		size INTEGER  NOT NULL
	);

	DROP TABLE IF EXISTS dir;
	CREATE TABLE dir (
		hash       CHAR(64) NOT NULL PRIMARY KEY,
		size       INTEGER  NOT NULL,
		count      INTEGER  NOT NULL,
		totalCount INTEGER  NOT NULL
	);

	DROP TABLE IF EXISTS dirItem;
	CREATE TABLE dirItem (
		hash           CHAR(64)  NOT NULL,
		itemHash       CHAR(64)  NOT NULL,
		itemName       TEXT      NOT NULL,
		itemMode       INTEGER   NOT NULL,
		itemSize       INTEGER   NOT NULL,
		itemCount      INTEGER   NOT NULL,
		itemTotalCount INTEGER   NOT NULL,
		itemCreateTime TIMESTAMP NOT NULL,
		itemModifyTime TIMESTAMP NOT NULL,
		itemChangeTime TIMESTAMP NOT NULL,
		itemAccessTime TIMESTAMP NOT NULL,
		PRIMARY KEY(Hash, itemName)
	);

	DROP TABLE IF EXISTS [commit];
	CREATE TABLE [commit] (
		id          INTEGER   NOT NULL PRIMARY KEY AUTOINCREMENT,
		createTime  TIMESTAMP NOT NULL,
		Hash        CHAR(64)  NOT NULL,
		lastId      INTEGER   NOT NULL
	);

	DROP TABLE IF EXISTS branch;
	CREATE TABLE branch (
		name        TEXT     NOT NULL PRIMARY KEY,
		description TEXT     NOT NULL DEFAULT "",
		commitId    INTEGER  NOT NULL,
		size        INTEGER  NOT NULL,
		count       INTEGER  NOT NULL
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
