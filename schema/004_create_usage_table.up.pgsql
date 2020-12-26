CREATE TABLE IF NOT EXISTS usages
(
    usage_id    BIGSERIAL NOT NULL UNIQUE,
    shortcut_id UUID      NOT NULL,
    logged_at   TIMESTAMP NOT NULL DEFAULT NOW(),
    PRIMARY KEY (usage_id),
    CONSTRAINT fk_shortcut_id FOREIGN KEY (shortcut_id) REFERENCES shortcuts (shortcut_id) ON DELETE CASCADE
);