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
SELECT * FROM _sync_from_baidu_photo;

-- name: InsertSyncFromBaiduPhoto :exec
INSERT INTO _sync_from_baidu_photo (
    id, 
    cron, 
    sync_at_start, 
    access_token, 
    refresh_token
)
VALUES (?, ?, ?, ?, ?);
-- name: UpsertSyncFromBaiduPhoto :exec
INSERT INTO _sync_from_baidu_photo (
    id, 
    cron, 
    sync_at_start, 
    access_token, 
    refresh_token
)
VALUES (?, ?, ?, ?, ?)
ON DUPLICATE KEY UPDATE
    cron = VALUES(cron),
    sync_at_start = VALUES(sync_at_start),
    access_token = VALUES(access_token),
    refresh_token = VALUES(refresh_token);
-- name: UpdateSyncFromBaiduPhoto :exec
UPDATE _sync_from_baidu_photo
SET 
    access_token = ?,
    refresh_token = ?
WHERE id = ?;
