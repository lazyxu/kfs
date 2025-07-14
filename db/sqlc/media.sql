-- FileType
-- name: InsertFileType :exec
INSERT INTO _file_type (hash, Type, SubType, Extension) VALUES (?, ?, ?, ?);
-- name: GetFileType :one
SELECT * FROM _file_type WHERE hash = ?;
-- name: ListFileType :many
SELECT * FROM _file_type;

-- DCIM Metadata
-- name: InsertDCIMMetadataTime :exec
INSERT INTO _dcim_metadata_time (hash, time, year, month, day) VALUES (?, ?, ?, ?, ?);
-- name: GetDCIMMetadataTime :one
SELECT * FROM _dcim_metadata_time WHERE hash = ?;

-- LivePhoto
-- name: InsertLivePhoto :exec
INSERT INTO _live_photo (movHash, heicHash, jpgHash, livpHash) VALUES (?, ?, ?, ?);
-- name: GetLivePhoto :one
SELECT * FROM _live_photo WHERE movHash = ?;
