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
		_logic_file.hash,
		_logic_file.parent_id,
		_logic_file.name,
		_file_type.sub_type,
		SUBSTRING_INDEX(_logic_file.name, '.', 1) AS perfix,
		CASE WHEN _logic_file.name LIKE '%.%' THEN LOWER(SUBSTRING_INDEX(_logic_file.name, '.', -1)) ELSE '' END AS suffix
	FROM _file_type
	INNER JOIN _logic_file ON _file_type.id = _logic_file.hash
	WHERE LOWER(SUBSTRING_INDEX(_logic_file.name, '.', -1)) = 'mov' AND _file_type.Extension = 'mov'
) AS t0
LEFT JOIN (
	SELECT
		_logic_file.hash,
		_logic_file.parent_id,
		_logic_file.name,
		SUBSTRING_INDEX(_logic_file.name, '.', 1) AS perfix,
		CASE WHEN _logic_file.name LIKE '%.%' THEN LOWER(SUBSTRING_INDEX(_logic_file.name, '.', -1)) ELSE '' END AS suffix,
		CASE WHEN LOWER(SUBSTRING_INDEX(_logic_file.name, '.', -1)) = 'heic' THEN _logic_file.hash ELSE NULL END AS heicHash,
		CASE WHEN LOWER(SUBSTRING_INDEX(_logic_file.name, '.', -1)) = 'jpg' THEN _logic_file.hash ELSE NULL END AS jpgHash
	FROM _file_type
	INNER JOIN _logic_file ON _file_type.id = _logic_file.hash
	WHERE (suffix = 'heic' AND _file_type.Extension = 'heif')
		OR (suffix = 'jpg' AND _file_type.Extension = 'jpg')
) AS t1
ON t0.parent_id = t1.parent_id AND t0.perfix = t1.perfix
ON DUPLICATE KEY UPDATE
	mov = VALUES(mov),
	heic = VALUES(heic),
	jpg = VALUES(jpg);
