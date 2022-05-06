package noncgo

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

type SqliteNonCgoDB struct {
	_db *sql.DB
}

func Open(dataSourceName string) (*SqliteNonCgoDB, error) {
	db, err := sql.Open("sqlite", dataSourceName)
	return &SqliteNonCgoDB{db}, err
}

func (db *SqliteNonCgoDB) Reset() error {
	_, err := db._db.Exec(`
	DROP TABLE IF EXISTS file;
	CREATE TABLE file (
		hash CHAR(64) NOT NULL PRIMARY KEY,
		size INTEGER  NOT NULL,
		ext  TEXT     NOT NULL
	);

	DROP TABLE IF EXISTS dir;
	CREATE TABLE dir (
		hash      CHAR(64) NOT NULL PRIMARY KEY,
		size      INTEGER  NOT NULL,
		count INTEGER  NOT NULL
	);

	DROP TABLE IF EXISTS dirItem;
	CREATE TABLE dirItem (
		hash           CHAR(64)  NOT NULL,
		itemHash       CHAR(64)  NOT NULL,
		itemName       TEXT      NOT NULL,
		itemMode       INTEGER   NOT NULL,
		itemSize       INTEGER   NOT NULL,
		itemCount      INTEGER   NOT NULL,
		itemCreateTime TIMESTAMP NOT NULL,
		itemModifyTime TIMESTAMP NOT NULL,
		itemChangeTime TIMESTAMP NOT NULL,
		itemAccessTime TIMESTAMP NOT NULL,
		oldItemHash    CHAR(64)  NOT NULL DEFAULT "",
		PRIMARY KEY(hash, itemName)
	);

	DROP TABLE IF EXISTS branch;
	CREATE TABLE branch (
		name        TEXT     NOT NULL PRIMARY KEY,
		description TEXT     NOT NULL DEFAULT "",
		hash        CHAR(64) NOT NULL,
		size        INTEGER  NOT NULL,
		count       INTEGER  NOT NULL
	);
	`)
	return err
}

func (db *SqliteNonCgoDB) Close() error {
	return db._db.Close()
}
