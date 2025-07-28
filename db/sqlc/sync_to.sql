-- name: ListLocalFileDriver :many
SELECT 
    *
FROM _local_file_sync_to_logic_file lfst
WHERE 
    lfst.device_id = ?;
