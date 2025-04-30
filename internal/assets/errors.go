package assets

import "github.com/rotisserie/eris"

var (
	// ErrValueMissing is returned when a value is missing.
	ErrValueMissing = eris.Errorf("missing value")

	// ErrValueInvalid is returned when the slug is invalid.
	ErrValueInvalid = eris.Errorf("invalid value")
)
