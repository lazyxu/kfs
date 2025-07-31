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
    sha256 CHAR(64) NOT NULL,
    UNIQUE INDEX idx_size_hash (size, md5, sha256)
);

CREATE TABLE IF NOT EXISTS _logic_file (
    id            BIGINT          NOT NULL PRIMARY KEY AUTO_INCREMENT,
    user_id       BIGINT          NOT NULL,
    parent_id     BIGINT          DEFAULT NULL,
    is_root       BOOLEAN GENERATED ALWAYS AS (IF(parent_id IS NULL, 1, NULL)) VIRTUAL, -- 新增虚拟列
    -- windows: 不允许使用的字符 \ / : * ? " < > |
    --          大小不敏感
    --          文件和文件夹名称不得超过255个字符，包括扩展名
    --          完整路径长度，最大为 260 个字符，启用了长路径支持，最大为 32,767 个字符
    --          避免使用 CON, PRN, AUX, NUL, COM1-COM9, LPT1-LPT9 等保留名称
    -- linux:   不允许使用的字符 /
    --          大小写敏感
    -- macos：  不允许使用的字符 :
    --          大小写敏感
    name          VARCHAR(255)    NOT NULL,
    -- current_version_id BIGINT         NOT NULL,
    file_id       BIGINT          DEFAULT NULL,
    mode          BIGINT          NOT NULL,  -- isDirectory, isFile, isSymlink
    size          BIGINT          NOT NULL,
    create_time   DATETIME        NOT NULL,  -- 创建时间，跨盘跨设备移动文件时也会更新，linux未强制要求记录创建时间
    modify_time   DATETIME        NOT NULL,  -- 内容修改时间，ls -l的时间
    change_time   DATETIME        NOT NULL,  -- 属性修改时间，windows不支持
    access_time   DATETIME        NOT NULL,  -- 访问时间
    FOREIGN KEY (parent_id)       REFERENCES _logic_file(id),
    FOREIGN KEY (file_id)         REFERENCES _file(id),
    -- FOREIGN KEY (current_version_id)  REFERENCES _version_file(id)
    UNIQUE INDEX idx_unique_user_id_parent_id_name (user_id, parent_id, name),
    UNIQUE INDEX idx_root_unique (user_id, is_root, name)
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
    type      VARCHAR(255) NOT NULL,
    sub_type  VARCHAR(255) NOT NULL,
    extension VARCHAR(255) NOT NULL,
    FOREIGN KEY (id)  REFERENCES   _file(id)
);

CREATE TABLE IF NOT EXISTS _height_width (
    id     BIGINT NOT NULL PRIMARY KEY,
    height BIGINT DEFAULT NULL,
    width  BIGINT DEFAULT NULL,
    FOREIGN KEY (id)  REFERENCES   _file(id)
);

CREATE TABLE IF NOT EXISTS _exif (
    id                    BIGINT        NOT NULL PRIMARY KEY,
    exif_version          CHAR(4)       DEFAULT NULL COMMENT 'Exif版本号',
    image_description     VARCHAR(255)  DEFAULT NULL COMMENT '图像描述',
    orientation           BIGINT        DEFAULT NULL COMMENT '方向（0~8）',
    date_time             VARCHAR(19)   DEFAULT NULL COMMENT '修改时间（YYYY:MM:DD HH:MM:SS）',
    date_time_original    VARCHAR(19)   DEFAULT NULL COMMENT '原始拍摄时间',
    date_time_digitized   VARCHAR(19)   DEFAULT NULL COMMENT '数字化时间',
    offset_time           VARCHAR(6)    DEFAULT NULL COMMENT '时区偏移（±HHMM）',
    offset_time_original  VARCHAR(6)    DEFAULT NULL COMMENT '原始时区偏移',
    offset_time_digitized VARCHAR(6)   DEFAULT NULL COMMENT '数字化时区偏移',
    subsec_time           VARCHAR(9)    DEFAULT NULL COMMENT '秒以下时间（毫秒/微秒）',
    subsec_time_original  VARCHAR(9)    DEFAULT NULL COMMENT '原始秒以下时间',
    subsec_time_digitized VARCHAR(9)   DEFAULT NULL COMMENT '数字化秒以下时间',
    host_computer         VARCHAR(255)  DEFAULT NULL COMMENT '处理设备型号',
    make                  VARCHAR(255)  DEFAULT NULL COMMENT '相机品牌',
    model                 VARCHAR(255)  DEFAULT NULL COMMENT '相机型号',
    exif_image_width      BIGINT        DEFAULT NULL COMMENT '图像宽度（像素）',
    exif_image_length     BIGINT        DEFAULT NULL COMMENT '图像高度（像素）',
    gps_latitude_ref      CHAR(1)       DEFAULT NULL COMMENT '纬度参考（N/S）',
    gps_latitude          DOUBLE        DEFAULT NULL COMMENT '纬度（十进制）',
    gps_longitude_ref     CHAR(1)       DEFAULT NULL COMMENT '经度参考（E/W）',
    gps_longitude         DOUBLE        DEFAULT NULL COMMENT '经度（十进制）',
    FOREIGN KEY (id) REFERENCES _file(id)
) COMMENT '存储图像Exif元数据';

-- TODO: heic -> jpg 只有文件头的uuid不一致
CREATE TABLE IF NOT EXISTS _content_same_file (
    id       BIGINT       NOT NULL PRIMARY KEY,
    id2      BIGINT       NOT NULL,
    FOREIGN KEY (id)   REFERENCES   _file(id),
    FOREIGN KEY (id2)  REFERENCES   _file(id)
);

CREATE TABLE IF NOT EXISTS _video_metadata (
    id       BIGINT       NOT NULL PRIMARY KEY,
    Codec    VARCHAR(255) DEFAULT NULL,
    Created  DATETIME     DEFAULT NULL,
    Modified DATETIME     DEFAULT NULL,
    Duration DOUBLE       DEFAULT NULL,
    FOREIGN KEY (id)  REFERENCES   _file(id)
);

CREATE TABLE IF NOT EXISTS _dcim_metadata_time (
    id       BIGINT     NOT NULL PRIMARY KEY,
    time     DATETIME   NOT NULL,
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

-- DROP TABLE IF EXISTS _sync_from_baidu_photo;

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
    cron          VARCHAR(64) NOT NULL,
    enabled       BOOLEAN NOT NULL DEFAULT FALSE,
    sync_at_start BOOLEAN NOT NULL DEFAULT FALSE,
    -- baidu_photo
    bduss_bfess   VARCHAR(192) NOT NULL COMMENT 'BDUSS_BFESS',
    stoken        VARCHAR(64)  NOT NULL COMMENT 'STOKEN',
    -- access_token  VARCHAR(256) NOT NULL,
    -- refresh_token VARCHAR(256) NOT NULL,
    FOREIGN KEY (id) REFERENCES _logic_file(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS _device (
    id         CHAR(36) NOT NULL PRIMARY KEY,
    name       VARCHAR(256) NOT NULL,
    os         VARCHAR(256) NOT NULL,
    user_agent VARCHAR(256) NOT NULL,
    host_name  VARCHAR(256) NOT NULL
);

-- DROP TABLE IF EXISTS _local_file_sync_to_logic_file;

CREATE TABLE IF NOT EXISTS _local_file_sync_to_logic_file (
    id            BIGINT        NOT NULL PRIMARY KEY,
    cron          VARCHAR(64)   NOT NULL,
    enabled       BOOLEAN       NOT NULL DEFAULT FALSE,
    sync_at_start BOOLEAN       NOT NULL DEFAULT FALSE,
    -- local_file
    device_id     CHAR(36)      NOT NULL,
    srcPath       VARCHAR(1024) NOT NULL,
    ignores       TEXT          NOT NULL,
    encoder       VARCHAR(64)   NOT NULL DEFAULT '',
    FOREIGN KEY (id)        REFERENCES _logic_file(id),
    FOREIGN KEY (device_id) REFERENCES _device(id)
);
