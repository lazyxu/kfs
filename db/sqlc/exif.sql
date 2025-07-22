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
-- name: UpsertNullExif :exec
INSERT INTO _exif (id) 
VALUES (?) 
ON DUPLICATE KEY UPDATE 
    exif_version = NULL,
    image_description = NULL,
    orientation = NULL,
    date_time = NULL,
    date_time_original = NULL,
    date_time_digitized = NULL,
    offset_time = NULL,
    offset_time_original = NULL,
    offset_time_digitized = NULL,
    subsec_time = NULL,
    subsec_time_original = NULL,
    subsec_time_digitized = NULL,
    host_computer = NULL,
    make = NULL,
    model = NULL,
    exif_image_width = NULL,
    exif_image_length = NULL,
    gps_latitude_ref = NULL,
    gps_latitude = NULL,
    gps_longitude_ref = NULL,
    gps_longitude = NULL;
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
-- name: UpsertVideoMetadata :exec
INSERT INTO _video_metadata (
    id,
    Codec, -- not null
    Created, -- not null
    Modified, -- not null
    Duration -- not null
) VALUES (?, ?, ?, ?, ?)
ON DUPLICATE KEY UPDATE
    Codec = VALUES(Codec),
    Created = VALUES(Created),
    Modified = VALUES(Modified),
    Duration = VALUES(Duration);
-- name: UpsertNullVideoMetadata :exec
INSERT INTO _video_metadata (id) 
VALUES (?)
ON DUPLICATE KEY UPDATE
    Codec = NULL,
    Created = NULL,
    Modified = NULL,
    Duration = NULL;
-- name: GetVideoMetadata :one
SELECT * FROM _video_metadata WHERE id = ?;
