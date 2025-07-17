-- name: InsertFile :execresult
INSERT INTO _file (size, md5, sha256) VALUES (?, ?, ?) ON DUPLICATE KEY UPDATE id=LAST_INSERT_ID(id);
-- name: GetFileByID :one
SELECT * FROM _file WHERE id = ?;
-- name: GetFileByHash :one
SELECT * FROM _file WHERE md5 = ? AND sha256 = ?;
-- name: ListFile :many
SELECT * FROM _file;
-- name: DeleteFileByID :exec
DELETE FROM _file WHERE id = ?;
-- name: SumFileSize :one
SELECT SUM(size) as size FROM _file;
-- name: GetMetadata :one
SELECT
    ft.type, ft.sub_type, ft.extension,
    ex.exif_version,
    ex.image_description,
    ex.orientation,
    ex.date_time,
    ex.date_time_original,
    ex.date_time_digitized,
    ex.offset_time,
    ex.offset_time_original,
    ex.offset_time_digitized,
    ex.subsec_time,
    ex.subsec_time_original,
    ex.subsec_time_digitized,
    ex.host_computer,
    ex.make,
    ex.model,
    ex.exif_image_width,
    ex.exif_image_length,
    ex.gps_latitude_ref,
    ex.gps_latitude,
    ex.gps_longitude_ref,
    ex.gps_longitude
FROM _file_type ft
LEFT JOIN _exif ex
    ON ft.id = ex.id  -- 关联条件保持不变
WHERE ft.id = ?;
-- name: ListFileExpectFileType :many
SELECT * FROM _file EXCEPT SELECT id FROM _file_type;
