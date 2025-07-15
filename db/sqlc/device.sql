-- Device
-- name: InsertDevice :exec
INSERT INTO _device (id, name, os, user_agent, host_name) VALUES (sqlc.arg('id'), sqlc.arg('name'), sqlc.arg('os'), sqlc.arg('user_agent'), sqlc.arg('host_name'))
ON DUPLICATE KEY UPDATE
    name=sqlc.arg('name'),
    os=sqlc.arg('os'),
    user_agent=sqlc.arg('user_agent'),
    host_name=sqlc.arg('host_name');
-- name: GetDevice :one
SELECT * FROM _device WHERE id = ?;
-- name: ListDevice :many
SELECT * FROM _device;
-- name: DeleteDevice :exec
DELETE FROM _device WHERE id = ?;
