-- name: InsertSyncFrom :execresult
INSERT INTO _sync_from_baidu_photo (
  id,
  cron
) VALUES (?, ?);

-- name: UpdateSyncFrom :exec
UPDATE _sync_from_baidu_photo
SET cron = ?
WHERE id = ?;

-- name: DeleteSyncFrom :exec
DELETE FROM _sync_from_baidu_photo
WHERE id = ?;

-- name: GetSyncFromById :one
SELECT * FROM _sync_from_baidu_photo
WHERE id = ?;

-- name: ListSync :many
SELECT * FROM _sync_from_baidu_photo WHERE sync = 1;
