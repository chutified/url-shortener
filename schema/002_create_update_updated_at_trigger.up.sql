CREATE OR REPLACE FUNCTION update_updated_at()
  RETURNS TRIGGER
  LANGUAGE PLPGSQL
  AS
$$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$;

CREATE TRIGGER update_updated_at_shortcuts
  BEFORE UPDATE
  ON shortcuts
  FOR EACH ROW
    EXECUTE PROCEDURE update_updated_at();
