package data

import "errors"

// ErrUnauthorized is returned if provided admin_key is invalid.
var ErrUnauthorized = errors.New("admin key validation failure")

// AdminAuth validates given admin key. ErrUnauthorized is returned
// if key is wrong. Otherwise unexpected internal server error is returned.
func (s *service) AdminAuth(key string) error {
	// TODO implement admin auth
	return nil
}
