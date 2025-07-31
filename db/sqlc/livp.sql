-- name: SetLivpForMovAndHeicOrJpgInDriver :execresult
INSERT INTO _live_photo (
	mov,
	heic,
	jpg
)
SELECT
	t0.hash AS mov,
	t1.heicHash,
	t1.jpgHash
FROM (
	SELECT
		lf.hash,
		lf.parent_id,
		lf.name,
		ft.sub_type,
		SUBSTRING_INDEX(lf.name, '.', 1) AS perfix,
		CASE WHEN lf.name LIKE '%.%' THEN LOWER(SUBSTRING_INDEX(lf.name, '.', -1)) ELSE '' END AS suffix
	FROM _file_type
	INNER JOIN _logic_file ON ft.id = lf.hash
	WHERE LOWER(SUBSTRING_INDEX(lf.name, '.', -1)) = 'mov' AND ft.Extension = 'mov'
) AS t0
LEFT JOIN (
	SELECT
		lf.hash,
		lf.parent_id,
		lf.name,
		SUBSTRING_INDEX(lf.name, '.', 1) AS perfix,
		CASE WHEN lf.name LIKE '%.%' THEN LOWER(SUBSTRING_INDEX(lf.name, '.', -1)) ELSE '' END AS suffix,
		CASE WHEN LOWER(SUBSTRING_INDEX(lf.name, '.', -1)) = 'heic' THEN lf.hash ELSE NULL END AS heicHash,
		CASE WHEN LOWER(SUBSTRING_INDEX(lf.name, '.', -1)) = 'jpg' THEN lf.hash ELSE NULL END AS jpgHash
	FROM _file_type
	INNER JOIN _logic_file ON ft.id = lf.hash
	WHERE (suffix = 'heic' AND ft.Extension = 'heif')
		OR (suffix = 'jpg' AND ft.Extension = 'jpg')
) AS t1
ON t0.parent_id = t1.parent_id AND t0.perfix = t1.perfix
ON DUPLICATE KEY UPDATE
	mov = VALUES(mov),
	heic = VALUES(heic),
	jpg = VALUES(jpg);

-- name: SetLivpForMovAndHeicOrJpgByParentId :exec
INSERT INTO _live_photo (mov, heic, jpg)
SELECT
    MAX(CASE WHEN suffix = 'mov' THEN file_id END) AS mov_id,
    MAX(CASE WHEN suffix = 'heic' THEN file_id END) AS heic_id,
    MAX(CASE WHEN suffix = 'jpg' THEN file_id END) AS jpg_id
FROM (
    SELECT 
        lf.parent_id,
        CASE 
            WHEN name LIKE '%.%' 
            THEN SUBSTRING(name, 1, LENGTH(name) - LOCATE('.', REVERSE(name))) 
            ELSE name 
        END AS prefix,
        LOWER(SUBSTRING_INDEX(name, '.', -1)) AS suffix,
        lf.file_id
    FROM _logic_file lf
    -- 添加文件类型验证
    INNER JOIN _file_type ft ON lf.file_id = ft.id  -- 假设_file_type.id关联_logic_file.file_id
    WHERE 
        ((lf.parent_id = sqlc.arg('parent_id')) OR (lf.parent_id IS NULL AND sqlc.arg('parent_id') IS NULL))
        AND (
            (LOWER(SUBSTRING_INDEX(lf.name, '.', -1)) = 'mov' AND ft.Extension = 'mov') OR
            (LOWER(SUBSTRING_INDEX(lf.name, '.', -1)) = 'heic' AND ft.Extension = 'heif') OR  -- HEIC扩展名在_file_type中可能为heif
            (LOWER(SUBSTRING_INDEX(lf.name, '.', -1)) = 'jpg' AND ft.Extension = 'jpg')
        )
        AND lf.file_id IS NOT NULL
) AS ValidFiles
GROUP BY parent_id, prefix
HAVING 
    mov_id IS NOT NULL 
    AND (heic_id IS NOT NULL OR jpg_id IS NOT NULL)
ON DUPLICATE KEY UPDATE
    heic = VALUES(heic),
    jpg = VALUES(jpg);

-- INSERT INTO _live_photo (
-- 	mov,
-- 	heic,
-- 	jpg
-- )
-- SELECT
-- 	t0.hash AS mov,
-- 	t1.heicHash,
-- 	t1.jpgHash
-- FROM (
-- 	SELECT
-- 		lf.hash,
-- 		lf.parent_id,
-- 		lf.name,
-- 		ft.sub_type,
-- 		SUBSTRING_INDEX(lf.name, '.', 1) AS perfix,
-- 		CASE WHEN lf.name LIKE '%.%' THEN LOWER(SUBSTRING_INDEX(lf.name, '.', -1)) ELSE '' END AS suffix
-- 	FROM _file_type ft
-- 	INNER JOIN _logic_file lf ON ft.id = lf.hash
-- 	WHERE ((lf.parent_id = sqlc.arg('parent_id'))
-- 			OR (lf.parent_id IS NULL AND sqlc.arg('parent_id') IS NULL))
-- 		AND (LOWER(SUBSTRING_INDEX(lf.name, '.', -1)) = 'mov' AND ft.Extension = 'mov')
-- ) AS t0
-- LEFT JOIN (
-- 	SELECT
-- 		lf.hash,
-- 		lf.parent_id,
-- 		lf.name,
-- 		SUBSTRING_INDEX(lf.name, '.', 1) AS perfix,
-- 		CASE WHEN lf.name LIKE '%.%' THEN LOWER(SUBSTRING_INDEX(lf.name, '.', -1)) ELSE '' END AS suffix,
-- 		CASE WHEN LOWER(SUBSTRING_INDEX(lf.name, '.', -1)) = 'heic' THEN lf.hash ELSE NULL END AS heicHash,
-- 		CASE WHEN LOWER(SUBSTRING_INDEX(lf.name, '.', -1)) = 'jpg' THEN lf.hash ELSE NULL END AS jpgHash
-- 	FROM _file_type ft
-- 	INNER JOIN _logic_file lf ON ft.id = lf.hash
-- 	WHERE ((lf.parent_id = sqlc.arg('parent_id'))
-- 			OR (lf.parent_id IS NULL AND sqlc.arg('parent_id') IS NULL))
-- 		AND ((suffix = 'heic' AND ft.Extension = 'heif')
-- 			OR (suffix = 'jpg' AND ft.Extension = 'jpg'))
-- ) AS t1
-- ON t0.parent_id = t1.parent_id AND t0.perfix = t1.perfix
-- ON DUPLICATE KEY UPDATE
-- 	mov = VALUES(mov),
-- 	heic = VALUES(heic),
-- 	jpg = VALUES(jpg);
