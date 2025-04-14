package tigris

import "errors"

var (
	// ErrMissingCreds is returned when the credentials are missing.
	ErrMissingCreds = errors.New("missing credentials (ENV)")

	// ErrInvalidCreds is returned when the credentials are invalid.
	ErrInvalidCreds = errors.New("invalid credentials (ENV)")
)
