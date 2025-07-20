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

-- name: UpsertLogicFile :execresult
INSERT INTO _logic_file (
    user_id,
    parent_id,
    name,
    file_id,
    mode,
    size,
    create_time,
    modify_time,
    change_time,
    access_time
)
VALUES (
    ?,  -- user_id
    ?,  -- parent_id
    ?,  -- name
    ?,  -- file_id
    ?,  -- mode
    ?,  -- size
    ?,  -- create_time
    ?,  -- modify_time
    ?,  -- change_time
    ?   -- access_time
)
ON DUPLICATE KEY UPDATE
    id = LAST_INSERT_ID(id),
    user_id = VALUES(user_id),
    file_id = VALUES(file_id),
    mode = VALUES(mode),
    size = VALUES(size),
    create_time = VALUES(create_time),
    modify_time = VALUES(modify_time),
    change_time = VALUES(change_time),
    access_time = VALUES(access_time);

-- name: DeleteLogicFile :exec
DELETE FROM _logic_file
WHERE id = ?;

-- name: GetLogicFileById :one
SELECT * FROM _logic_file
WHERE id = ?;

-- name: GetLogicFile :one
SELECT * FROM _logic_file
WHERE user_id = ? AND (
  (parent_id = sqlc.arg('parent_id'))
  OR (parent_id IS NULL AND sqlc.arg('parent_id') IS NULL)
) AND name = ?;

-- name: ListLogicFilesByUserId :many
SELECT * FROM _logic_file
WHERE user_id = ?;

-- name: ListLogicFilesByParentId :many
SELECT
  lf.id,
  lf.user_id,
  lf.parent_id,
  lf.is_root,
  lf.name,
  lf.file_id,
  lf.mode,
  lf.size,
  lf.create_time,
  lf.modify_time,
  lf.change_time,
  lf.access_time,
  COALESCE(f.md5, '') AS md5,
  COALESCE(f.sha256, '') AS sha256
FROM _logic_file lf
LEFT JOIN _file f ON lf.file_id = f.id  
WHERE user_id = ? AND (
  (parent_id = sqlc.arg('parent_id'))
  OR (parent_id IS NULL AND sqlc.arg('parent_id') IS NULL)
);

-- name: GetEarliestCrated :one
SELECT LEAST(
    create_time,
    modify_time,
    change_time,
    access_time
) AS min_time
FROM _logic_file
WHERE file_id = ?;
