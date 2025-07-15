-- LogicFile

-- name: InsertLogicFile :execresult
INSERT INTO _logic_file (
  parent_id,
  user_id,
  name,
  file_id,
  mode,
  size,
  create_time,
  modify_time,
  change_time,
  access_time
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?);

-- name: UpdateLogicFile :exec
UPDATE _logic_file
SET parent_id = ?, user_id = ?, name = ?, file_id = ?, mode = ?, size = ?, create_time = ?, modify_time = ?, change_time = ?, access_time = ?
WHERE id = ?;

-- name: DeleteLogicFile :exec
DELETE FROM _logic_file
WHERE id = ?;

-- name: GetLogicFileById :one
SELECT * FROM _logic_file
WHERE id = ?;

-- name: GetLogicFile :one
SELECT * FROM _logic_file
WHERE user_id = ? AND (
  (parent_id = ?)
  OR (parent_id IS NULL AND ? IS NULL)
) AND name = ?;

-- 说明：
-- MySQL 不支持 IS NOT DISTINCT FROM（仅 PostgreSQL 和 MySQL 8.0.22+ 支持）。
-- 推荐用上述 OR (parent_id IS NULL AND ? IS NULL) 方式，兼容所有 MySQL 版本。

-- name: ListLogicFileByUserId :many
SELECT * FROM _logic_file
WHERE user_id = ?;

-- name: ListLogicFileByParentId :many
SELECT * FROM _logic_file
WHERE parent_id = ?;
