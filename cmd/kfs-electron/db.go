package main

import (
	"context"
	"database/sql"
	_ "modernc.org/sqlite"
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
	`)
	return err
}

func (db *DB) Create() error {
	conn := db.getConn()
	defer db.putConn(conn)
	_, err := conn.Exec(`
	CREATE TABLE IF NOT EXISTS _file (
	    time       INTEGER      NOT NULL,
		path       VARCHAR(256) NOT NULL,
	    dirname    VARCHAR(256) NOT NULL,
		name       VARCHAR(256) NOT NULL,
	    typ        INTEGER      NOT NULL,
		count      INTEGER      NOT NULL,
		size       INTEGER      NOT NULL,
	    PRIMARY KEY(time, path)
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

func (db *DB) InsertFile(ctx context.Context, time int64, path string, isDir bool, count int64, size int64) error {
	conn := db.getConn()
	defer db.putConn(conn)
	_, err := conn.ExecContext(ctx, `
	INSERT INTO _file VALUES (?, ?, ?, ?, ?, ?, ?);
	`, time, path, filepath.Dir(path), filepath.Base(path), isDir, count, size)
	return err
}
