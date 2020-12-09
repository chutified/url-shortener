-- name: IncrementUsage :execresult
UPDATE
  shortcuts
SET
  usage = usage + 1
WHERE
  shortcut_id = $1
  AND deleted_at = NULL
RETURNING
  shortcut_id,
  usage;

-- name: GetTotalUsage :one
SELECT
  SUM(usage)
FROM
  shotcuts;
