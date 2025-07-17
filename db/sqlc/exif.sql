-- Exif
-- name: InsertExif :exec
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
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);
-- name: InsertNullExif :exec
INSERT INTO _exif ( id ) VALUES (?);
-- name: GetExif :one
SELECT * FROM _exif WHERE id = ?;
-- name: ListExif :many
SELECT * FROM _exif;

-- HeightWidth
-- name: InsertHeightWidth :exec
INSERT INTO _height_width (id, height, width) VALUES (?, ?, ?);
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
