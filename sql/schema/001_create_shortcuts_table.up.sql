-- +migrate Up
CREATE TABLE shortcuts (
  shortcut_id UUID NOT NULL UNIQUE,
  full_url TEXT NOT NULL,
  short_url VARCHAR(255) NOT NULL UNIQUE,
  usage INTEGER NOT NULL,
  create_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  deleted_at TIMESTAMP,
  PRIMARY KEY (shortcut_id)
);
