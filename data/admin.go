package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/chutommy/rand"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

//goland:noinspection ALL
const (
	salt = "@salt"

	prefixLen = 8
	keyLen    = 40
	digits    = "0123456789"
	alphabet  = "ABCDEFGHIJKLMNOPQRSTUVWXYZ" + "abcdefghijklmnopqrstuvwxyz"
	charSet   = alphabet + digits

	hashedPasswd = "$2a$10$SBWWLZ4QvaTeUNk1moBW9O29Vuf4/KiXPweTcakYm4X1onaS/ZA1m" //nolint:gosec
	username     = "urlshorteneradmin"
)

var keySplitLen = 2

var (
	// ErrUnauthorized is returned if provided admin_key is invalid.
	ErrUnauthorized = errors.New("admin key validation failure")
	// ErrPrefixNotFound is returned if admin_key with the given prefix can not be found.
	ErrPrefixNotFound = errors.New("admin_key's prefix was not found")
)

// AuthenticateAdmin validates if the given passwd is correct.
func (s *service) AuthenticateAdmin(name string, passwd string) error {
	if name != username {
		return ErrUnauthorized
	}

	// validate
	err := bcrypt.CompareHashAndPassword([]byte(hashedPasswd), []byte(passwd+salt))
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return ErrUnauthorized
	} else if err != nil {
		return fmt.Errorf("unexpected validation failure: %w", err)
	}

	// success
	return nil
}

// ValidateAdminKey validates given admin key. ErrUnauthorized is returned
// if key is wrong. Otherwise, unexpected internal server error is returned.
func (s *service) ValidateAdminKey(ctx context.Context, wholeKey string) error {
	wholeKey += salt

	// separate the wholeKey
	splitKey := strings.Split(wholeKey, ".")
	if len(splitKey) != keySplitLen {
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
	if err := row.Scan(&hashKey); errors.Is(err, sql.ErrNoRows) {
		return ErrUnauthorized
	} else if err != nil {
		return fmt.Errorf("retrieved sql row scan error: %w", err)
	}

	// compare
	if err := bcrypt.CompareHashAndPassword(
		[]byte(hashKey),
		[]byte(key),
	); errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return ErrUnauthorized
	} else if err != nil {
		return fmt.Errorf("unexpected validation failure: %w", err)
	}

	return nil
}

// GenerateAdminKey generates a new admin_key and add it into the database.
func (s *service) GenerateAdminKey(ctx context.Context) (string, error) {
	var prefix, key []byte

	for {
		// generate key
		prefix, key = genKey()

		// hash key
		hashKey, err := bcrypt.GenerateFromPassword(append(key, []byte(salt)...), bcrypt.DefaultCost)
		if err != nil {
			return "", fmt.Errorf("unable to hash generated password: %w", err)
		}

		// insert
		_, err = s.DB.ExecContext(ctx, `
INSERT INTO
  admin_keys (prefix, hashed_key)
VALUES
  ($1, $2);
  `, prefix, hashKey)
		if err != nil {
			// postgres errors
			var pqErr *pq.Error
			if errors.As(err, pqErr) {
				// unique violation
				if pqErr.Code == "23505" {
					continue
				}
			}

			return "", fmt.Errorf("insert failure: %w", err)
		}

		break
	}

	return string(prefix) + "." + string(key), nil
}

// RevokeAdminKey revokes admin_key with the given unique prefix.
func (s *service) RevokeAdminKey(ctx context.Context, prefix string) error {
	// revoke
	res, err := s.DB.ExecContext(ctx, `
UPDATE
  admin_keys
SET
  revoked_at = NOW()
WHERE
  prefix = $1
  AND revoked_at IS NULL;
  `, prefix)
	if err != nil {
		return fmt.Errorf("unexpected validation failure: %w", err)
	}

	// check result
	if i, _ := res.RowsAffected(); i == 0 {
		return ErrPrefixNotFound
	}

	return nil
}

// genKey generates a random prefix and key of an api key.
func genKey() ([]byte, []byte) {
	// init buffers
	prefix := make([]byte, prefixLen)
	key := make([]byte, keyLen)

	r := rand.New()
	// fill
	for i := 0; i < prefixLen; i++ {
		prefix[i] = alphabet[r.Intn(len(alphabet))]
	}

	for i := 0; i < keyLen; i++ {
		key[i] = charSet[r.Intn(len(charSet))]
	}

	// shuffle
	r.Shuffle(prefixLen, func(i, j int) {
		prefix[i], prefix[j] = prefix[j], prefix[i]
	})
	r.Shuffle(keyLen, func(i, j int) {
		key[i], key[j] = key[j], key[i]
	})

	return prefix, key
}
