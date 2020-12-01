-- name: GetRecordsLen :one
SELECT
  COUNT(*)
FROM
  shortcuts
WHERE
  deleted_at = NULL;

-- name: GetAllRecords :many
SELECT
  shortcut_id,
  full_url,
  short_url
FROM
  shortcuts
WHERE
  deleted_at = NULL
ORDER BY
  $1
LIMIT
  $2 OFFSET $3;

-- name: GetDetails :one
SELECT
  shortcut_id,
  full_url,
  short_url,
  usage,
  created_at,
  updated_at
FROM
  shortcuts
WHERE
  shortcut_id = $1
  AND deleted_at = NULL
LIMIT
  1;

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
