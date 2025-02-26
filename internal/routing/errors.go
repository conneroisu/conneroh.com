package routing

import (
	"fmt"
	"net/url"
)

// ErrNotFound is an error that is returned when a resource is not found.
type ErrNotFound struct {
	URL *url.URL
}

// Error implements the error interface on ErrNotFound.
func (e ErrNotFound) Error() string {
	return fmt.Sprintf(
		"resource not found for %s",
		e.URL.String(),
	)
}

// ErrMissingParam is an error that is returned when a required parameter is
// missing.
type ErrMissingParam struct {
	ID   string
	View string
}

// Error implements the error interface on ErrMissingParam.
func (e ErrMissingParam) Error() string {
	return fmt.Sprintf(
		"missing parameter for %s: %s",
		e.View,
		e.ID,
	)
}
