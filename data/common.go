package data

import (
	"context"
	"database/sql"
	"fmt"
)

// newNullString returns a passed string in sql's NullString type.
func newNullString(s string) sql.NullString {
	if len(s) == 0 {
		return sql.NullString{}
	}

	return sql.NullString{
		String: s,
		Valid:  true,
	}
}

// LogError logs the error into error_logs table.
func (s *service) LogError(ctx context.Context, logErr error) {
	// insert
	_, err := s.DB.ExecContext(ctx, `
INSERT INTO
  error_logs (err_msg)
VALUES
  ($1);
	`, logErr)
	if err != nil {
		fmt.Printf("[ERROR] unalbe to log error (%v): %v", logErr, err)
	}
}
