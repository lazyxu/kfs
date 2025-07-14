-- DriverFile
-- name: InsertDriverFile :exec
INSERT INTO _driver_file (driverId, dirPath, name, hash, mode, size, createTime, modifyTime, changeTime, accessTime) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?);
-- name: GetDriverFile :one
SELECT * FROM _driver_file WHERE driverId = ? AND dirPath = ? AND name = ?;
-- name: ListDriverFile :many
SELECT * FROM _driver_file WHERE driverId = ?;

-- DriverFileHistory
-- name: InsertDriverFileHistory :exec
INSERT INTO _driver_file_history (driverId, dirPath, name, hash, mode, size, createTime, modifyTime, changeTime, accessTime, uploadDeviceId, uploadTime) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);
-- name: GetDriverFileHistory :one
SELECT * FROM _driver_file_history WHERE driverId = ? AND dirPath = ? AND name = ? AND uploadDeviceId = ? AND uploadTime = ?;
