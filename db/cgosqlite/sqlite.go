package cgosqlite

import (
	"database/sql"

	"github.com/lazyxu/kfs/dao"
	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	dataSourceName string
	ch             chan *sql.DB
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
	conn, err := sql.Open("sqlite3", dataSourceName)
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

func (db *DB) Remove() error {
	conn := db.getConn()
	defer db.putConn(conn)
	_, err := conn.Exec(`
	DROP TABLE IF EXISTS _file;
	DROP TABLE IF EXISTS _dir;
	DROP TABLE IF EXISTS _dirItem;
	DROP TABLE IF EXISTS _commit;
	DROP TABLE IF EXISTS _branch;
	`)
	return err
}

func (db *DB) Create() error {
	conn := db.getConn()
	defer db.putConn(conn)
	_, err := conn.Exec(`
	CREATE TABLE IF NOT EXISTS _file (
		hash CHAR(64) NOT NULL PRIMARY KEY,
		size INTEGER  NOT NULL
	);

	CREATE TABLE IF NOT EXISTS _dir (
		hash       CHAR(64) NOT NULL PRIMARY KEY,
		size       INTEGER  NOT NULL,
		count      INTEGER  NOT NULL,
		totalCount INTEGER  NOT NULL
	);

	CREATE TABLE IF NOT EXISTS _dirItem (
		hash           CHAR(64)     NOT NULL,
		itemHash       CHAR(64)     NOT NULL,
		itemName       VARCHAR(256) NOT NULL,
		itemMode       INTEGER       NOT NULL,
		itemSize       INTEGER       NOT NULL,
		itemCount      INTEGER       NOT NULL,
		itemTotalCount INTEGER       NOT NULL,
		itemCreateTime TIMESTAMP    NOT NULL,
		itemModifyTime TIMESTAMP    NOT NULL,
		itemChangeTime TIMESTAMP    NOT NULL,
		itemAccessTime TIMESTAMP    NOT NULL,
		PRIMARY KEY(Hash, itemName)
	);

	CREATE TABLE IF NOT EXISTS _commit (
		id          INTEGER    NOT NULL PRIMARY KEY AUTOINCREMENT,
		createTime  TIMESTAMP NOT NULL,
		Hash        CHAR(64)  NOT NULL,
		lastId      INTEGER    NOT NULL
	);

	CREATE TABLE IF NOT EXISTS _branch (
		name        VARCHAR(256) NOT NULL PRIMARY KEY,
		description VARCHAR(256) NOT NULL DEFAULT "",
		commitId    INTEGER       NOT NULL,
		size        INTEGER       NOT NULL,
		count       INTEGER       NOT NULL,
		FOREIGN KEY (commitId)   REFERENCES _commit(id)
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
