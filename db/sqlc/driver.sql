-- Driver
-- name: InsertDriver :exec
INSERT INTO _driver (name, description, Type) VALUES (?, ?, ?);
-- name: GetDriver :one
SELECT * FROM _driver WHERE id = ?;
-- name: ListDriver :many
SELECT * FROM _driver;

-- DriverSync
-- name: InsertDriverSync :exec
INSERT INTO _driver_sync (id, sync, h, m) VALUES (?, ?, ?, ?);
-- name: GetDriverSync :one
SELECT * FROM _driver_sync WHERE id = ?;

-- DriverBaiduPhoto
-- name: InsertDriverBaiduPhoto :exec
INSERT INTO _driver_baidu_photo (id, accessToken, refreshToken) VALUES (?, ?, ?);
-- name: GetDriverBaiduPhoto :one
SELECT * FROM _driver_baidu_photo WHERE id = ?;

-- DriverLocalFile
-- name: InsertDriverLocalFile :exec
INSERT INTO _driver_local_file (id, deviceId, srcPath, ignores, encoder) VALUES (?, ?, ?, ?, ?);
-- name: GetDriverLocalFile :one
SELECT * FROM _driver_local_file WHERE id = ?;
