CREATE TABLE IF NOT EXISTS shortcuts (
  shortcut_id UUID PRIMARY KEY,
  full_url TEXT NOT NULL,
  short_url VARCHAR(255) NOT NULL UNIQUE,
  usage INTEGER NOT NULL,
  create_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  deleted_at TIMESTAMP
);

-- name: AddRecord :execresult
INSERT INTO
  shortcuts (shortcut_id, full_url, short_url)
VALUES
  ($1, $2, $3);

-- name: UpdateRecord :execresult
UPDATE
  shortcuts
SET
  full_url = COALESCE($2, full_url),
  short_url = COALESCE($3, short_url)
WHERE
  shortcut_id = $1
  AND deleted_at = NULL
RETURNING
  shortcut_id,
  full_url,
  short_url;

-- name: DeleteRecord :execresult
UPDATE
  shortcuts
SET
  deleted_at = NOW()
WHERE
  shortcut_id = $1
  AND deleted_at = NULL
RETURNING
  shortcut_id,
  full_url,
  short_url;

-- name: GetRecordByID :one
SELECT
  shortcut_id,
  full_url,
  short_url
FROM
  shortcuts
WHERE
  shortcut_id = $1
  AND deleted_at = NULL
LIMIT
  1;

-- name: GetRecordByShort :one
SELECT
  shortcut_id,
  full_url,
  short_url
FROM
  shortcuts
WHERE
  short_url = $1
  AND deleted_at = NULL
LIMIT
  1;

-- name: GetRecordByFull :one
SELECT
  shortcut_id,
  full_url,
  short_url
FROM
  shortcuts
WHERE
  full_url = $1
  AND deleted_at = NULL
LIMIT
  1;

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

-- name: RecordRecovery :execresult
UPDATE
  shortcuts
SET
  deleted_at = NULL
WHERE
  deleted_at != NULL
  AND shortcut_id = $1
RETURNING
  shortcut_id,
  full_url,
  short_url;

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
