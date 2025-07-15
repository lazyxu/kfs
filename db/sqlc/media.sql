-- FileType
-- name: InsertFileType :exec
INSERT INTO _file_type (id, Type, SubType, Extension) VALUES (?, ?, ?, ?);
-- name: GetFileType :one
SELECT * FROM _file_type WHERE id = ?;
-- name: ListFileType :many
SELECT * FROM _file_type;

-- DCIM Metadata
-- name: InsertDCIMMetadataTime :exec
INSERT INTO _dcim_metadata_time (id, time, year, month, day) VALUES (?, ?, ?, ?, ?);
-- name: GetDCIMMetadataTime :one
SELECT * FROM _dcim_metadata_time WHERE id = ?;

-- LivePhoto
-- name: InsertLivePhoto :exec
INSERT INTO _live_photo (mov, heic, jpg, livp) VALUES (?, ?, ?, ?);
-- name: GetLivePhoto :one
SELECT * FROM _live_photo WHERE mov = ?;
