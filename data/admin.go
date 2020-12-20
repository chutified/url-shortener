package data

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

var (
	salt = "@salt"
	// ErrUnauthorized is returned if provided admin_key is invalid.
	ErrUnauthorized = errors.New("admin key validation failure")
)

// AdminAuth validates given admin key. ErrUnauthorized is returned
// if key is wrong. Otherwise unexpected internal server error is returned.
func (s *service) AdminAuth(ctx context.Context, wholeKey string) error {

	wholeKey = strings.ToLower(wholeKey) + salt

	// seperate the wholeKey
	splitKey := strings.Split(wholeKey, ".")
	if len(splitKey) != 2 {
		return ErrUnauthorized
	}

	prefix := splitKey[0]
	key := splitKey[1]

	// query db
	row := s.DB.QueryRowxContext(ctx, `
SELECT
  hashed_key
FROM
  admin_keys
WHERE
  prefix = $1
  AND revoked_at IS NULL;
	`, prefix)

	// scan row
	var hashKey string
	if err := row.Scan(&hashKey); err == sql.ErrNoRows {
		return ErrUnauthorized
	} else if err != nil {
		return errors.New("unexpected internal server error")
	}

	// compare
	if err := bcrypt.CompareHashAndPassword([]byte(hashKey), []byte(key)); err == bcrypt.ErrMismatchedHashAndPassword {
		return ErrUnauthorized
	} else if err != nil {
		return errors.New("unexpected error")
	}

	return nil
}
