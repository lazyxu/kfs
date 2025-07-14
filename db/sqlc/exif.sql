-- Exif
-- name: InsertExif :exec
INSERT INTO _exif (hash, ExifVersion, ImageDescription, Orientation, DateTime, DateTimeOriginal, DateTimeDigitized, OffsetTime, OffsetTimeOriginal, OffsetTimeDigitized, SubsecTime, SubsecTimeOriginal, SubsecTimeDigitized, HostComputer, Make, Model, ExifImageWidth, ExifImageLength, GPSLatitudeRef, GPSLatitude, GPSLongitudeRef, GPSLongitude) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);
-- name: GetExif :one
SELECT * FROM _exif WHERE hash = ?;
-- name: ListExif :many
SELECT * FROM _exif;

-- HeightWidth
-- name: InsertHeightWidth :exec
INSERT INTO _height_width (hash, height, width) VALUES (?, ?, ?);
-- name: GetHeightWidth :one
SELECT * FROM _height_width WHERE hash = ?;

-- VideoMetadata
-- name: InsertVideoMetadata :exec
INSERT INTO _video_metadata (hash, Codec, Created, Modified, Duration) VALUES (?, ?, ?, ?, ?);
-- name: GetVideoMetadata :one
SELECT * FROM _video_metadata WHERE hash = ?;
