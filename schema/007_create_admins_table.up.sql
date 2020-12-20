CREATE TABLE IF NOT EXISTS admins (
  admin_id UUID NOT NULL UNIQUE,
  username VARCHAR(255) NOT NULL UNIQUE,
  hashed_passwd VARCHAR(32) NOT NULL UNIQUE,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  PRIMARY KEY (admin_id)
);

CREATE OR REPLACE FUNCTION set_created_at_admins()
  RETURNS TRIGGER
  LANGUAGE PLPGSQL
  AS $$
BEGIN
  NEW.created_at = NOW();
  RETURN NEW;
END;
$$;

CREATE TRIGGER set_created_at_admins_trigger
  BEFORE INSERT
  ON admins
  FOR EACH ROW
    EXECUTE PROCEDURE set_created_at_admins();
