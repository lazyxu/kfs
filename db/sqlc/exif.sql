-- Exif
-- name: InsertExif :exec
INSERT INTO _exif (id, ExifVersion, ImageDescription, Orientation, DateTime, DateTimeOriginal, DateTimeDigitized, OffsetTime, OffsetTimeOriginal, OffsetTimeDigitized, SubsecTime, SubsecTimeOriginal, SubsecTimeDigitized, HostComputer, Make, Model, ExifImageWidth, ExifImageLength, GPSLatitudeRef, GPSLatitude, GPSLongitudeRef, GPSLongitude) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);
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
INSERT INTO _video_metadata (id, Codec, Created, Modified, Duration) VALUES (?, ?, ?, ?, ?);
-- name: GetVideoMetadata :one
SELECT * FROM _video_metadata WHERE id = ?;
