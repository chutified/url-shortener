CREATE TABLE IF NOT EXISTS shortcuts (
  shortcut_id UUID PRIMARY KEY,
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
  ($1, $2, $3);

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
