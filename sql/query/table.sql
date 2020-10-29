CREATE TABLE IF NOT EXISTS shortcuts (
  shortcut_id UUID PRIMARY KEY,
  full_url TEXT NOT NULL,
  short_url VARCHAR(255) NOT NULL UNIQUE,
  usage INTEGER NOT NULL,
  create_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  deleted_at TIMESTAMP
);
