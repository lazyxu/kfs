
CREATE TABLE IF NOT EXISTS _file (
    hash CHAR(64) NOT NULL PRIMARY KEY,
    size BIGINT   NOT NULL
);

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
    ExifImageWidth      BIGINT        DEFAULT NULL,
    ExifImageLength     BIGINT        DEFAULT NULL,
    GPSLatitudeRef      CHAR(1)      DEFAULT NULL,
    GPSLatitude         DOUBLE       DEFAULT NULL,
    GPSLongitudeRef     CHAR(1)      DEFAULT NULL,
    GPSLongitude        DOUBLE       DEFAULT NULL,
    FOREIGN KEY (hash)  REFERENCES   _file(hash)
);

CREATE TABLE IF NOT EXISTS _driver (
    id BIGINT NOT NULL PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(256) NOT NULL UNIQUE,
    description VARCHAR(256) NOT NULL DEFAULT '',
    Type VARCHAR(256) NOT NULL DEFAULT ''
);

CREATE TABLE IF NOT EXISTS _driver_file (
    driverId    BIGINT        NOT NULL,
    dirPath    VARCHAR(376) NOT NULL,
    name       VARCHAR(255)   NOT NULL,
    version    BIGINT         NOT NULL,
    hash       CHAR(64)       NOT NULL,
    mode       BIGINT         NOT NULL,
    size       BIGINT         NOT NULL,
    createTime BIGINT         NOT NULL,
    modifyTime BIGINT         NOT NULL,
    changeTime BIGINT         NOT NULL,
    accessTime BIGINT         NOT NULL,
    PRIMARY KEY (driverId, dirPath, name, version),
    FOREIGN KEY (driverId)  REFERENCES _driver(id)
);

-- 其它表名（字段需补充）：

CREATE TABLE IF NOT EXISTS _file_md5 (
    hash CHAR(64) NOT NULL PRIMARY KEY,
    md5  CHAR(32) NOT NULL
);
-- CREATE INDEX index_md5 ON _file_md5 (md5);

CREATE TABLE IF NOT EXISTS _height_width (
    hash CHAR(64) NOT NULL PRIMARY KEY,
    height BIGINT DEFAULT NULL,
    width BIGINT DEFAULT NULL,
    FOREIGN KEY (hash) REFERENCES _file(hash)
);

CREATE TABLE IF NOT EXISTS _video_metadata (
    hash CHAR(64) NOT NULL PRIMARY KEY,
    Codec VARCHAR(255) DEFAULT NULL,
    Created BIGINT DEFAULT NULL,
    Modified BIGINT DEFAULT NULL,
    Duration DOUBLE DEFAULT NULL,
    FOREIGN KEY (hash) REFERENCES _file(hash)
);

CREATE TABLE IF NOT EXISTS _file_type (
    hash CHAR(64) NOT NULL PRIMARY KEY,
    Type VARCHAR(255) NOT NULL,
    SubType VARCHAR(255) NOT NULL,
    Extension VARCHAR(255) NOT NULL,
    FOREIGN KEY (hash) REFERENCES _file(hash)
);

CREATE TABLE IF NOT EXISTS _dcim_metadata_time (
    hash CHAR(64) NOT NULL PRIMARY KEY,
    time BIGINT DEFAULT NULL,
    year BIGINT DEFAULT NULL,
    month BIGINT DEFAULT NULL,
    day BIGINT DEFAULT NULL,
    FOREIGN KEY (hash) REFERENCES _file(hash)
);

CREATE TABLE IF NOT EXISTS _live_photo (
    movHash CHAR(64) NOT NULL PRIMARY KEY,
    heicHash CHAR(64) DEFAULT NULL,
    jpgHash CHAR(64) DEFAULT NULL,
    livpHash CHAR(64) DEFAULT NULL,
    FOREIGN KEY (movHash) REFERENCES _file(hash),
    FOREIGN KEY (heicHash) REFERENCES _file(hash),
    FOREIGN KEY (jpgHash) REFERENCES _file(hash),
    FOREIGN KEY (livpHash) REFERENCES _file(hash)
);

CREATE TABLE IF NOT EXISTS _driver_sync (
    id BIGINT NOT NULL PRIMARY KEY,
    sync TINYINT NOT NULL DEFAULT 1,
    h TINYINT NOT NULL DEFAULT 4,
    m TINYINT NOT NULL DEFAULT 0,
    FOREIGN KEY (id) REFERENCES _driver(id)
);

CREATE TABLE IF NOT EXISTS _driver_baidu_photo (
    id BIGINT NOT NULL PRIMARY KEY,
    accessToken VARCHAR(256) NOT NULL DEFAULT '',
    refreshToken VARCHAR(256) NOT NULL DEFAULT '',
    FOREIGN KEY (id) REFERENCES _driver(id)
);

CREATE TABLE IF NOT EXISTS _device (
    id CHAR(36) NOT NULL PRIMARY KEY,
    name VARCHAR(256) NOT NULL,
    os VARCHAR(256) NOT NULL,
    userAgent VARCHAR(256) NOT NULL,
    hostname VARCHAR(256) NOT NULL
);

CREATE TABLE IF NOT EXISTS _driver_local_file (
    id BIGINT NOT NULL PRIMARY KEY,
    deviceId CHAR(36) NOT NULL DEFAULT '',
    srcPath VARCHAR(1024) NOT NULL DEFAULT '',
    ignores TEXT NOT NULL,
    encoder VARCHAR(64) NOT NULL DEFAULT '',
    FOREIGN KEY (id) REFERENCES _driver(id),
    FOREIGN KEY (deviceId) REFERENCES _device(id)
);

CREATE TABLE IF NOT EXISTS _driver_file_history (
    driverId BIGINT NOT NULL,
    dirPath VARCHAR(376) NOT NULL,
    name VARCHAR(255) NOT NULL,
    hash CHAR(64) NOT NULL,
    mode BIGINT NOT NULL,
    size BIGINT NOT NULL,
    createTime BIGINT NOT NULL,
    modifyTime BIGINT NOT NULL,
    changeTime BIGINT NOT NULL,
    accessTime BIGINT NOT NULL,
    uploadDeviceId CHAR(36) NOT NULL,
    uploadTime BIGINT NOT NULL,
    PRIMARY KEY (driverId, dirPath, name, uploadDeviceId, uploadTime),
    FOREIGN KEY (driverId) REFERENCES _driver(id),
    FOREIGN KEY (uploadDeviceId) REFERENCES _device(id)
);
