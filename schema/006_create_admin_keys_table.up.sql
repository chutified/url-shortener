CREATE TABLE IF NOT EXISTS admin_keys (
  key_id BIGSERIAL NOT NULL UNIQUE,
  prefix VARCHAR(8) NOT NULL UNIQUE,
  hashed_key VARCHAR(32) NOT NULL UNIQUE,
  generated_at TIMESTAMP NOT NULL DEFAULT NOW(),
  revoked_at TIMESTAMP DEFAULT NULL,
  PRIMARY KEY (key_id)
);

CREATE OR REPLACE FUNCTION set_generated_at_admin_keys()
  RETURNS TRIGGER
  LANGUAGE PLPGSQL
  AS
$$
BEGIN
  NEW.generated_at = NOW();
  NEW.revoked_at = NULL;
RETURN NEW;
END;
$$;

CREATE TRIGGER set_generated_at_admin_keys_trigger BEFORE
INSERT
  ON admin_keys
  FOR EACH ROW
    EXECUTE PROCEDURE set_generated_at_admin_keys();
