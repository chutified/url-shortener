package data

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"errors"
	"strings"
)

// ErrUnauthorized is returned if provided admin_key is invalid.
var ErrUnauthorized = errors.New("admin key validation failure")

// AdminAuth validates given admin key. ErrUnauthorized is returned
// if key is wrong. Otherwise unexpected internal server error is returned.
func (s *service) AdminAuth(ctx context.Context, key string) error {

	key = strings.ToLower(key)

	// seperate the key
	splitKey := strings.Split(key, ".")
	prefix := splitKey[0]
	hashKey := sha256.Sum256([]byte(splitKey[1]))

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
	var dbKey [32]byte
	if err := row.Scan(&dbKey); err == sql.ErrNoRows {
		return ErrUnauthorized
	} else if err != nil {
		return errors.New("unexpected internal server error")
	}

	// compare
	if hashKey != dbKey {
		return ErrUnauthorized
	}

	return nil
}
