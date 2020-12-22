package data

import (
	"context"
	"database/sql"
	"errors"
	"math/rand"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const (
	salt = "@salt"

	// generating admin_key options
	prefixLen = 8
	keyLen    = 40
	digits    = "0123456789"
	specials  = "?!@#$%^&*()[]{}<>_=-+|;:"
	charSet   = "ABCDEFGHIJKLMNOPQRSTUVWXYZ" + "abcdefghijklmnopqrstuvwxyz" + digits + specials
)

var (
	hashedPasswd = []byte("$2a$10$SBWWLZ4QvaTeUNk1moBW9O29Vuf4/KiXPweTcakYm4X1onaS/ZA1m")
	username     = "urlshorteneradmin"

	// ErrUnauthorized is returned if provided admin_key is invalid.
	ErrUnauthorized = errors.New("admin key validation failure")
)

// AuthenticateAdmin validates if the given passwd is correct.
func (s *service) AuthenticateAdmin(ctx context.Context, name string, passwd string) error {

	// check username
	if name != username {
		return ErrUnauthorized
	}

	// validate
	err := bcrypt.CompareHashAndPassword(hashedPasswd, []byte(passwd+salt))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return ErrUnauthorized
	} else if err != nil {
		return errors.New("unexpected error when comparing hashed password")
	}

	// success
	return nil
}

// GenerateAdminKey generates a new admin_key and add it into the database.
func (s *service) GenerateAdminKey(ctx context.Context) (string, error) {

	// set seed
	rand.Seed(time.Now().UnixNano())

	// init buffers
	prefix := make([]byte, prefixLen)
	key := make([]byte, keyLen)

	// add digit
	prefix[0] = digits[rand.Intn(len(digits))]
	key[0] = digits[rand.Intn(len(digits))]

	// add special char
	prefix[1] = specials[rand.Intn(len(specials))]
	key[1] = specials[rand.Intn(len(specials))]

	// fill
	for i := 2; i < prefixLen; i++ {
		prefix[i] = charSet[rand.Intn(len(charSet))]
	}
	for i := 2; i < keyLen; i++ {
		key[i] = charSet[rand.Intn(len(charSet))]
	}

	// shuffle
	rand.Shuffle(prefixLen, func(i, j int) {
		prefix[i], prefix[j] = prefix[j], prefix[i]
	})
	rand.Shuffle(keyLen, func(i, j int) {
		key[i], key[j] = key[j], key[i]
	})

	// build
	newKey := string(prefix) + "." + string(key)

	// TODO insert into database

	return newKey, nil
}

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
