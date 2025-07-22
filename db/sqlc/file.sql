-- name: UpsertFile :execresult
INSERT INTO _file (size, md5, sha256) VALUES (?, ?, ?) ON DUPLICATE KEY UPDATE id=LAST_INSERT_ID(id);
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
-- name: ListFileExpectFileType :many
SELECT id, size, md5, sha256 FROM _file EXCEPT SELECT id, NULL, NULL, NULL FROM _file_type;
