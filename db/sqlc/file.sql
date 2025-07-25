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
-- name: ListFileNotAnalyzed :many
SELECT id, size, md5, sha256
FROM _file f
WHERE NOT EXISTS (
    SELECT 1 FROM _file_type ft 
    WHERE ft.id = f.id  -- 通过 id 关联排除
);