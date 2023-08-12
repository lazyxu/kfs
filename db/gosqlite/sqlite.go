package gosqlite

import (
	"database/sql"

	"github.com/lazyxu/kfs/dao"
	_ "modernc.org/sqlite"
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

	CREATE TABLE IF NOT EXISTS _driver (
		name        VARCHAR(256) NOT NULL PRIMARY KEY,
		description VARCHAR(256) NOT NULL DEFAULT ""
	);

	CREATE TABLE IF NOT EXISTS _driver_file (
		driver_name VARCHAR(256)   NOT NULL,
		dirpath     VARCHAR(32767) NOT NULL,
		name        VARCHAR(255)   NOT NULL,
	    version     INTEGER        NOT NULL,
		hash        CHAR(64)       NOT NULL,
		mode        INTEGER        NOT NULL,
		size        INTEGER        NOT NULL,
		createTime  INTEGER        NOT NULL,
		modifyTime  INTEGER        NOT NULL,
		changeTime  INTEGER        NOT NULL,
		accessTime  INTEGER        NOT NULL,
		PRIMARY KEY (driver_name, dirpath, name, version),
		FOREIGN KEY (driver_name)  REFERENCES _driver(name)
	);
	`) // https://blog.csdn.net/jimmyleeee/article/details/124682486
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
