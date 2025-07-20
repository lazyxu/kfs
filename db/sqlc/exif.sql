-- Exif
-- name: UpsertExif :exec
INSERT INTO _exif (
    id,
    exif_version,
    image_description,
    orientation,
    date_time,
    date_time_original,
    date_time_digitized,
    offset_time,
    offset_time_original,
    offset_time_digitized,
    subsec_time,
    subsec_time_original,
    subsec_time_digitized,
    host_computer,
    make,
    model,
    exif_image_width,
    exif_image_length,
    gps_latitude_ref,
    gps_latitude,
    gps_longitude_ref,
    gps_longitude
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
ON DUPLICATE KEY UPDATE
    exif_version = VALUES(exif_version),
    image_description = VALUES(image_description),
    orientation = VALUES(orientation),
    date_time = VALUES(date_time),
    date_time_original = VALUES(date_time_original),
    date_time_digitized = VALUES(date_time_digitized),
    offset_time = VALUES(offset_time),
    offset_time_original = VALUES(offset_time_original),
    offset_time_digitized = VALUES(offset_time_digitized),
    subsec_time = VALUES(subsec_time),
    subsec_time_original = VALUES(subsec_time_original),
    subsec_time_digitized = VALUES(subsec_time_digitized),
    host_computer = VALUES(host_computer),
    make = VALUES(make),
    model = VALUES(model),
    exif_image_width = VALUES(exif_image_width),
    exif_image_length = VALUES(exif_image_length),
    gps_latitude_ref = VALUES(gps_latitude_ref),
    gps_latitude = VALUES(gps_latitude),
    gps_longitude_ref = VALUES(gps_longitude_ref),
    gps_longitude = VALUES(gps_longitude);
-- name: InsertNullExif :exec
INSERT INTO _exif ( id ) VALUES (?);
-- name: GetExif :one
SELECT * FROM _exif WHERE id = ?;
-- name: ListExif :many
SELECT * FROM _exif;

-- HeightWidth
-- name: UpsertHeightWidth :exec
INSERT INTO _height_width (id, height, width)
VALUES (?, ?, ?)
ON DUPLICATE KEY UPDATE
    height = VALUES(height),
    width = VALUES(width);
-- name: GetHeightWidth :one
SELECT * FROM _height_width WHERE id = ?;

-- VideoMetadata
-- name: InsertVideoMetadata :exec
INSERT INTO _video_metadata (
    id,
    Codec, -- not null
    Created, -- not null
    Modified, -- not null
    Duration -- not null
) VALUES (?, ?, ?, ?, ?);
-- name: InsertNullVideoMetadata :exec
INSERT INTO _video_metadata ( id ) VALUES (?);
-- name: GetVideoMetadata :one
SELECT * FROM _video_metadata WHERE id = ?;
