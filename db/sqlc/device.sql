-- Device
-- name: InsertDevice :exec
INSERT INTO _device (id, name, os, userAgent, hostname) VALUES (?, ?, ?, ?, ?);
-- name: GetDevice :one
SELECT * FROM _device WHERE id = ?;
-- name: ListDevice :many
SELECT * FROM _device;
