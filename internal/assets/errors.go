package assets

import (
	"errors"

	"github.com/rotisserie/eris"
)

var (
	// ErrValueMissing is returned when a value is missing.
	ErrValueMissing = eris.Errorf("missing value")

	// ErrValueInvalid is returned when the slug is invalid.
	ErrValueInvalid = eris.Errorf("invalid value")

	// ErrMissingCreds is returned when the credentials are missing.
	ErrMissingCreds = errors.New("missing credentials (ENV)")

	// ErrInvalidCreds is returned when the credentials are invalid.
	ErrInvalidCreds = errors.New("invalid credentials (ENV)")
)
