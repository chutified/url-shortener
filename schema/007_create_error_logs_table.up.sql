CREATE TABLE IF NOT EXISTS error_logs (
  err_id BIGSERIAL NOT NULL UNIQUE,
  err_msg TEXT NOT NULL,
  logged_at TIMESTAMP NOT NULL DEFAULT NOW(),
  PRIMARY KEY (err_id)
);

CREATE OR REPLACE FUNCTION set_logged_at ()
  RETURNS TRIGGER LANGUAGE plpgsql AS $$
BEGIN
  NEW.logged_at = NOW();
  RETURN NEW;
END;
$$;

CREATE TRIGGER set_logged_at_trigger
  BEFORE INSERT ON error_logs
  FOR EACH ROW EXECUTE PROCEDURE set_logged_at ();
