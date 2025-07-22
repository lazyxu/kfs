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
INSERT INTO _dcim_metadata_time (id, time)
VALUES (?, ?)
ON DUPLICATE KEY UPDATE
    time = VALUES(time);
-- name: GetDCIMMetadataTime :one
SELECT * FROM _dcim_metadata_time WHERE id = ?;

-- LivePhoto
-- name: InsertLivePhoto :exec
INSERT INTO _live_photo (mov, heic, jpg, livp) VALUES (?, ?, ?, ?);
-- name: GetLivePhoto :one
SELECT * FROM _live_photo WHERE mov = ?;

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

-- name: GetDCIM :one
SELECT
    f.id,
    f.size,
    f.md5,
    f.sha256,
    ft.type,
    ft.sub_type,
    ft.extension,
    hw.height,
    hw.width,
    dmt.time,
    vm.duration,
    e.exif_version,
    e.image_description,
    e.orientation,
    e.date_time,
    e.date_time_original,
    e.date_time_digitized,
    e.offset_time,
    e.offset_time_original,
    e.offset_time_digitized,
    e.subsec_time,
    e.subsec_time_original,
    e.subsec_time_digitized,
    e.host_computer,
    e.make,
    e.model,
    e.exif_image_width,
    e.exif_image_length,
    e.gps_latitude_ref,
    e.gps_latitude,
    e.gps_longitude_ref,
    e.gps_longitude
FROM _dcim_metadata_time dmt
INNER JOIN _file f ON dmt.id = f.id 
INNER JOIN _height_width hw ON dmt.id = hw.id  
INNER JOIN _file_type ft ON dmt.id = ft.id
LEFT JOIN _video_metadata vm ON dmt.id = vm.id
LEFT JOIN _exif e ON dmt.id = e.id
WHERE dmt.id = ?;
-- name: ListDCIM :many
SELECT
    f.id,
    f.size,
    f.md5,
    f.sha256,
    ft.type,
    ft.sub_type,
    ft.extension,
    hw.height,
    hw.width,
    dmt.time,
    vm.duration,
    e.exif_version,
    e.image_description,
    e.orientation,
    e.date_time,
    e.date_time_original,
    e.date_time_digitized,
    e.offset_time,
    e.offset_time_original,
    e.offset_time_digitized,
    e.subsec_time,
    e.subsec_time_original,
    e.subsec_time_digitized,
    e.host_computer,
    e.make,
    e.model,
    e.exif_image_width,
    e.exif_image_length,
    e.gps_latitude_ref,
    e.gps_latitude,
    e.gps_longitude_ref,
    e.gps_longitude
FROM _dcim_metadata_time dmt
INNER JOIN _file f ON dmt.id = f.id 
INNER JOIN _height_width hw ON dmt.id = hw.id  
INNER JOIN _file_type ft ON dmt.id = ft.id
LEFT JOIN _video_metadata vm ON dmt.id = vm.id
LEFT JOIN _exif e ON dmt.id = e.id
ORDER BY dmt.time DESC;
