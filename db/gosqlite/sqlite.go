package gosqlite

import (
	"database/sql"
	"os"

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
		ch:             make(chan *sql.DB, 1),
		dataSourceName: dataSourceName,
	}
	db.ch <- conn
	return db, err
}

func (db *DB) IsSqlite() bool {
	return true
}

func (db *DB) DataSourceName() string {
	return db.dataSourceName
}

func (db *DB) Size() (int64, error) {
	stat, err := os.Stat(db.dataSourceName)
	if err != nil {
		return 0, err
	}
	return stat.Size(), err
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
	DROP TABLE IF EXISTS _file_md5;
	DROP TABLE IF EXISTS _exif;
	DROP TABLE IF EXISTS _height_width;
	DROP TABLE IF EXISTS _video_metadata;
	DROP TABLE IF EXISTS _file_type;
	DROP TABLE IF EXISTS _live_photo;
	DROP TABLE IF EXISTS _driver;
	DROP TABLE IF EXISTS _driver_file;
	`)
	return err
}

func (db *DB) Create() error {
	conn := db.getConn()
	defer db.putConn(conn)
	// TODO: id as PRIMARY KEY
	// id   INT64    NOT NULL PRIMARY KEY AUTOINCREMENT,
	// hash CHAR(64) NOT NULL UNIQUE,
	_, err := conn.Exec(`
	CREATE TABLE IF NOT EXISTS _file (
		hash CHAR(64) NOT NULL PRIMARY KEY,
		size INT64    NOT NULL
	);
	CREATE TABLE IF NOT EXISTS _file_md5 (
		hash CHAR(64) NOT NULL PRIMARY KEY,
		md5  CHAR(32) NOT NULL
	);
	CREATE INDEX IF NOT EXISTS index_md5 ON _file_md5 (md5);

	CREATE TABLE IF NOT EXISTS _exif (
		hash                CHAR(64)     NOT NULL PRIMARY KEY,
		ExifVersion         CHAR(4)      DEFAULT NULL,
		ImageDescription    VARCHAR(255) DEFAULT NULL,
		Orientation         INT8         DEFAULT NULL,
		DateTime            VARCHAR(19)  DEFAULT NULL,
		DateTimeOriginal    VARCHAR(19)  DEFAULT NULL,
		DateTimeDigitized   VARCHAR(19)  DEFAULT NULL,
		OffsetTime          VARCHAR(6)   DEFAULT NULL,
		OffsetTimeOriginal  VARCHAR(6)   DEFAULT NULL,
		OffsetTimeDigitized VARCHAR(6)   DEFAULT NULL,
		SubsecTime          VARCHAR(9)   DEFAULT NULL,
		SubsecTimeOriginal  VARCHAR(9)   DEFAULT NULL,
		SubsecTimeDigitized VARCHAR(9)   DEFAULT NULL,
		HostComputer        VARCHAR(255) DEFAULT NULL,
		Make                VARCHAR(255) DEFAULT NULL,
		Model               VARCHAR(255) DEFAULT NULL,
		ExifImageWidth      INT64        DEFAULT NULL,
		ExifImageLength     INT64        DEFAULT NULL,
		GPSLatitudeRef      CHAR(1)      DEFAULT NULL,
		GPSLatitude         DOUBLE       DEFAULT NULL,
		GPSLongitudeRef     CHAR(1)      DEFAULT NULL,
		GPSLongitude        DOUBLE       DEFAULT NULL,
		FOREIGN KEY (hash)  REFERENCES   _file(hash)
	);

	CREATE TABLE IF NOT EXISTS _height_width (
		hash       CHAR(64)     NOT NULL PRIMARY KEY,
		height     INT64        DEFAULT NULL,
		width      INT64        DEFAULT NULL,
		FOREIGN KEY (hash)  REFERENCES   _file(hash)
	);

	CREATE TABLE IF NOT EXISTS _video_metadata (
		hash                CHAR(64)     NOT NULL PRIMARY KEY,
		Codec               VARCHAR(255) DEFAULT NULL,
		Created             INT64        DEFAULT NULL,
		Modified            INT64        DEFAULT NULL,
		Duration            DOUBLE       DEFAULT NULL,
		FOREIGN KEY (hash)  REFERENCES   _file(hash)
	);

	CREATE TABLE IF NOT EXISTS _file_type (
		hash               CHAR(64)     NOT NULL PRIMARY KEY,
		Type               VARCHAR(255) NOT NULL,
		SubType            VARCHAR(255) NOT NULL,
		Extension          VARCHAR(255) NOT NULL,
		FOREIGN KEY (hash) REFERENCES   _file(hash)
	);

	CREATE TABLE IF NOT EXISTS _live_photo (
		movHash                CHAR(64)     NOT NULL PRIMARY KEY,
		heicHash               CHAR(64)     DEFAULT NULL,
		jpgHash                CHAR(64)     DEFAULT NULL,
		FOREIGN KEY (movHash)  REFERENCES   _file(hash),
		FOREIGN KEY (heicHash) REFERENCES   _file(hash),
		FOREIGN KEY (jpgHash)  REFERENCES   _file(hash)
	);

	CREATE TABLE IF NOT EXISTS _driver (
		name         VARCHAR(256) NOT NULL PRIMARY KEY,
		description  VARCHAR(256) NOT NULL DEFAULT "",
		Type         VARCHAR(256) NOT NULL DEFAULT "",
		accessToken  VARCHAR(256) NOT NULL DEFAULT "",
		refreshToken VARCHAR(256) NOT NULL DEFAULT ""
	);

	CREATE TABLE IF NOT EXISTS _driver_file (
		driverName  VARCHAR(256)   NOT NULL,
		dirPath     VARCHAR(32767) NOT NULL,
		name        VARCHAR(255)   NOT NULL,
	    version     INT64          NOT NULL,
		hash        CHAR(64)       NOT NULL,
		mode        INT64          NOT NULL,
		size        INT64          NOT NULL,
		createTime  INT64          NOT NULL,
		modifyTime  INT64          NOT NULL,
		changeTime  INT64          NOT NULL,
		accessTime  INT64          NOT NULL,
		PRIMARY KEY (driverName, dirPath, name, version),
		FOREIGN KEY (driverName)  REFERENCES _driver(name)
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
