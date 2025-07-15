CREATE TABLE IF NOT EXISTS _user (
  id          BIGINT       NOT NULL PRIMARY KEY AUTO_INCREMENT,
  email       VARCHAR(255) NOT NULL UNIQUE,
  password    VARCHAR(255) NOT NULL,
  is_verified BOOLEAN      NOT NULL DEFAULT FALSE,
  is_deleted  BOOLEAN      NOT NULL DEFAULT FALSE,
  created_at  DATETIME     NOT NULL DEFAULT NOW(),
  updated_at  DATETIME     NOT NULL DEFAULT NOW() ON UPDATE NOW() 
);

CREATE TABLE IF NOT EXISTS _file (
    id     BIGINT   NOT NULL PRIMARY KEY AUTO_INCREMENT,
    size   BIGINT   NOT NULL,
    md5    CHAR(32) NOT NULL,
    sha256 CHAR(64) NOT NULL
);

CREATE TABLE IF NOT EXISTS _logic_file (
    id            BIGINT         NOT NULL PRIMARY KEY AUTO_INCREMENT,
    parent_id     BIGINT         DEFAULT NULL,
    user_id       BIGINT         NOT NULL,
    name          VARCHAR(255)   NOT NULL,
    -- current_version_id BIGINT         NOT NULL,
    file_id       BIGINT          DEFAULT NULL,
    mode          BIGINT          NOT NULL,  -- isDirectory, isFile, isSymlink
    size          BIGINT          NOT NULL,
    create_time   BIGINT          NOT NULL,
    modify_time   BIGINT          NOT NULL,
    change_time   BIGINT          NOT NULL,
    access_time   BIGINT          NOT NULL,
    -- FOREIGN KEY (parent_id)       REFERENCES _logic_file(id),
    FOREIGN KEY (file_id)         REFERENCES _file(id)
    -- FOREIGN KEY (current_version_id)  REFERENCES _version_file(id)
);

-- CREATE TABLE IF NOT EXISTS _version_file (
--     id            BIGINT          NOT NULL PRIMARY KEY AUTO_INCREMENT,
--     description   VARCHAR(255)    NOT NULL,
--     logic_file_id BIGINT          NOT NULL,
--     file_id       BIGINT          NOT NULL,
--     mode          BIGINT          NOT NULL,  -- isDirectory, isFile, isSymlink
--     size          BIGINT          NOT NULL,
--     create_time   BIGINT          NOT NULL,
--     modify_time   BIGINT          NOT NULL,
--     change_time   BIGINT          NOT NULL,
--     access_time   BIGINT          NOT NULL,
--     FOREIGN KEY (logic_file_id)   REFERENCES _logic_file(id),
--     FOREIGN KEY (file_id)         REFERENCES _file(id)
-- );

CREATE TABLE IF NOT EXISTS _file_type (
    id        BIGINT       NOT NULL PRIMARY KEY,
    Type      VARCHAR(255) NOT NULL,
    SubType   VARCHAR(255) NOT NULL,
    Extension VARCHAR(255) NOT NULL,
    FOREIGN KEY (id)  REFERENCES   _file(id)
);

CREATE TABLE IF NOT EXISTS _height_width (
    id     BIGINT NOT NULL PRIMARY KEY,
    height BIGINT DEFAULT NULL,
    width  BIGINT DEFAULT NULL,
    FOREIGN KEY (id)  REFERENCES   _file(id)
);

CREATE TABLE IF NOT EXISTS _exif (
    id                  BIGINT       NOT NULL PRIMARY KEY,
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
    FOREIGN KEY (id)  REFERENCES   _file(id)
);

CREATE TABLE IF NOT EXISTS _video_metadata (
    id       BIGINT       NOT NULL PRIMARY KEY,
    Codec    VARCHAR(255) DEFAULT NULL,
    Created  BIGINT DEFAULT NULL,
    Modified BIGINT DEFAULT NULL,
    Duration DOUBLE DEFAULT NULL,
    FOREIGN KEY (id)  REFERENCES   _file(id)
);

CREATE TABLE IF NOT EXISTS _dcim_metadata_time (
    id       BIGINT       NOT NULL PRIMARY KEY,
    time     BIGINT DEFAULT NULL,
    year     BIGINT DEFAULT NULL,
    month    BIGINT DEFAULT NULL,
    day      BIGINT DEFAULT NULL,
    FOREIGN KEY (id) REFERENCES _file(id)
);

CREATE TABLE IF NOT EXISTS _live_photo (
    mov  BIGINT NOT NULL PRIMARY KEY,
    heic BIGINT DEFAULT NULL,
    jpg  BIGINT DEFAULT NULL,
    livp BIGINT DEFAULT NULL,
    FOREIGN KEY (mov)  REFERENCES _file(id),
    FOREIGN KEY (heic) REFERENCES _file(id),
    FOREIGN KEY (jpg)  REFERENCES _file(id),
    FOREIGN KEY (livp) REFERENCES _file(id)
);

CREATE TABLE IF NOT EXISTS _sync_from_baidu_photo (
    id   BIGINT NOT NULL PRIMARY KEY,
    -- Cron任务的时间表达式由5个字段组成，分别表示：
    -- 分钟 小时 日期 月份 星期几

    -- 每个字段的取值范围及含义如下：
    -- 字段	取值范围	说明
    -- ​分钟​	0-59	每小时的第几分钟执行
    -- ​小时​	0-23	每天的第几小时执行
    -- ​日期​	1-31	每月的第几天执行
    -- ​月份​	1-12	每年的第几月执行
    -- ​星期几​	0-7	每周的第几天执行（0和7均代表周日）

    -- ⚙️ ​二、特殊符号详解​
    -- 以下符号用于灵活定义时间：
    -- 符号	示例	说明
    -- *	* * * * *	匹配所有值（如*在分钟字段=每分钟）
    -- ,	0,15,30 * * * *	指定多个时间点（每小时的0、15、30分执行）
    -- -	0 9-18 * * *	定义连续范围（每天9点至18点整点执行）
    -- */n	*/10 * * * *	固定间隔（每10分钟执行一次）

    -- ​三、常用时间设置示例​
    -- 表达式	说明
    -- * * * * *	​每分钟执行一次
    -- 0 * * * *	​每小时整点执行（如1:00, 2:00）
    -- 0 2 * * *	​每天凌晨2点执行
    -- 0 0 * * 0	​每周日0点执行
    -- 0 0 1 * *	​每月1日0点执行
    -- 0 9-18 * * 1-5	​工作日9点至18点整点每小时执行
    -- */5 * * * *	​每5分钟执行一次
    -- 0 20 * * 1,3,5	​每周一、三、五晚上8点执行
    cron VARCHAR(64) NOT NULL,
    sync BOOLEAN NOT NULL DEFAULT FALSE,
    -- baidu_photo
    access_token  VARCHAR(256) NOT NULL,
    refresh_token VARCHAR(256) NOT NULL,
    FOREIGN KEY (id) REFERENCES _logic_file(id)
);

CREATE TABLE IF NOT EXISTS _device (
    id         CHAR(36) NOT NULL PRIMARY KEY,
    name       VARCHAR(256) NOT NULL,
    os         VARCHAR(256) NOT NULL,
    user_agent VARCHAR(256) NOT NULL,
    host_name  VARCHAR(256) NOT NULL
);

CREATE TABLE IF NOT EXISTS _local_file_sync_to_logic_file (
    id        BIGINT        NOT NULL PRIMARY KEY,
    device_id CHAR(36)      NOT NULL DEFAULT '',
    srcPath   VARCHAR(1024) NOT NULL DEFAULT '',
    ignores   TEXT          NOT NULL,
    encoder   VARCHAR(64)   NOT NULL DEFAULT '',
    FOREIGN KEY (id)        REFERENCES _logic_file(id),
    FOREIGN KEY (device_id) REFERENCES _device(id)
);
