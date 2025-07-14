-- 文件相关
-- name: InsertFile :exec
INSERT INTO _file (hash, size) VALUES (?, ?);
-- name: GetFile :one
SELECT * FROM _file WHERE hash = ?;
-- name: ListFile :many
SELECT * FROM _file;
-- name: DeleteFile :exec
DELETE FROM _file WHERE hash = ?;
-- name: SumFileSize :one
SELECT SUM(size) as size FROM _file;

-- 文件MD5
-- name: InsertFileMd5 :exec
INSERT INTO _file_md5 (hash, md5) VALUES (?, ?);
-- name: GetFileByMd5 :one
SELECT * FROM _file_md5 WHERE md5 = ?;
