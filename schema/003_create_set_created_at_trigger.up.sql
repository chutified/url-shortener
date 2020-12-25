CREATE OR REPLACE FUNCTION set_init_timestamp()
  RETURNS TRIGGER LANGUAGE PLPGSQL AS $$
BEGIN
  NEW.created_at = NOW();
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$;

CREATE TRIGGER set_init_timestamp_shortcuts
  BEFORE INSERT ON shortcuts
  FOR EACH ROW
    EXECUTE PROCEDURE set_init_timestamp();
