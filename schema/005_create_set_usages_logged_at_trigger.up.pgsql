CREATE OR REPLACE FUNCTION set_logged_at_timestamp ()
    RETURNS TRIGGER
    LANGUAGE PLPGSQL
    AS $$
BEGIN
    NEW.logged_at = NOW();
    RETURN NEW;
END;
$$;

CREATE TRIGGER set_logged_at_timestamp_usages
    BEFORE INSERT ON usages
    FOR EACH ROW
    EXECUTE PROCEDURE set_logged_at_timestamp ();
