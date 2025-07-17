-- FileType
-- name: UpsertFileType :exec
INSERT INTO _file_type (id, type, sub_type, extension)
VALUES (?, ?, ?, ?)
ON DUPLICATE KEY UPDATE
    type = VALUES(type),
    sub_type = VALUES(sub_type),
    extension = VALUES(extension);
-- name: GetFileType :one
SELECT * FROM _file_type WHERE id = ?;
-- name: ListFileType :many
SELECT * FROM _file_type;

-- DCIM Metadata
-- name: UpsertDCIMMetadataTime :exec
INSERT INTO _dcim_metadata_time (id, time, year, month, day)
VALUES (?, ?, ?, ?, ?)
ON DUPLICATE KEY UPDATE
    time = VALUES(time),
    year = VALUES(year),
    month = VALUES(month),
    day = VALUES(day);
-- name: GetDCIMMetadataTime :one
SELECT * FROM _dcim_metadata_time WHERE id = ?;

-- LivePhoto
-- name: InsertLivePhoto :exec
INSERT INTO _live_photo (mov, heic, jpg, livp) VALUES (?, ?, ?, ?);
-- name: GetLivePhoto :one
SELECT * FROM _live_photo WHERE mov = ?;
