-- name: InsertFile :execresult
INSERT INTO _file (size, md5, sha256)
VALUES (?, ?, ?);
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
