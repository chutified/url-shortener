DROP TRIGGER IF EXISTS set_logged_at_trigger on error_logs;
DROP FUNCTION IF EXISTS set_logged_at();
DROP TABLE IF EXISTS error_logs;