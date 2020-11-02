CREATE TABLE IF NOT EXISTS shortcuts (
  shortcut_id UUID NOT NULL UNIQUE PRIMARY KEY,
  full_url TEXT NOT NULL,
  short_url VARCHAR(255) NOT NULL UNIQUE,
  usage INTEGER NOT NULL,
  create_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  deleted_at TIMESTAMP
);

-- AddRecord
INSERT INTO
  shortcuts (shortcut_id, full_url, short_url)
VALUES
  ($1, $2, $3)
ON CONFLICT (short_url) DO NOTHING;

-- UpdateRecord
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

-- DeleteRecord
UPDATE
  shortcuts
SET
  deleted_at = LOCALTIMESTAMP
WHERE
  shortcut_id = $1
  AND deleted_at = NULL
RETURNING
  shortcut_id,
  full_url,
  short_url;

-- GetRecordByID
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

-- GetRecordByShort
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

-- GetRecordByFull
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

-- GetRecordsLen
SELECT
  COUNT(*)
FROM
  shortcuts
WHERE
  deleted_at = NULL;

-- GetAllRecords
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

-- GetDetails
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

-- RecordRecovery
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

-- IncrementUsage
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

-- GetTotalUsage
SELECT
  SUM(usage)
FROM
  shotcuts;
