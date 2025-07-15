-- 文件相关
-- name: InsertFile :execresult
INSERT INTO _file (size, md5, sha256)
VALUES (?, ?, ?);
-- name: GetFile :one
SELECT * FROM _file WHERE id = ?;
-- name: ListFile :many
SELECT * FROM _file;
-- name: DeleteFile :exec
DELETE FROM _file WHERE id = ?;
-- name: SumFileSize :one
SELECT SUM(size) as size FROM _file;
