-- name: GetSyncFromById :one
SELECT * FROM _sync_from_baidu_photo
WHERE id = ?;

-- name: ListSync :many
SELECT * FROM _sync_from_baidu_photo;

-- name: InsertSyncFromBaiduPhoto :exec
INSERT INTO _sync_from_baidu_photo (
    id,
    cron,
    enabled,
    sync_at_start,
    bduss_bfess,
    stoken
)
VALUES (?, ?, ?, ?, ?, ?);

-- name: UpsertSyncFromBaiduPhoto :exec
INSERT INTO _sync_from_baidu_photo (
    id,
    cron,
    enabled,
    sync_at_start,
    bduss_bfess,
    stoken
)
VALUES (?, ?, ?, ?, ?, ?)
ON DUPLICATE KEY UPDATE
    cron = VALUES(cron),
    enabled = VALUES(enabled),
    sync_at_start = VALUES(sync_at_start),
    bduss_bfess = VALUES(bduss_bfess),
    stoken = VALUES(stoken);
