CREATE TABLE IF NOT EXISTS shortcuts
(
    shortcut_id UUID         NOT NULL UNIQUE,
    full_url    TEXT         NOT NULL,
    short_url   VARCHAR(255) NOT NULL UNIQUE,
    usage       INTEGER      NOT NULL DEFAULT 0,
    created_at  TIMESTAMP    NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMP    NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMP             DEFAULT NULL,
    PRIMARY KEY (shortcut_id)
);