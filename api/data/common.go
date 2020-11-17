package data

import "database/sql"

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
