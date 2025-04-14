package genstruct

import (
	"fmt"
	"reflect"
)

// NonSliceOrArrayError is returned when the data is not a slice or array.
type NonSliceOrArrayError struct {
	Kind reflect.Kind
}

// Error returns the error message
func (e NonSliceOrArrayError) Error() string {
	return fmt.Sprintf(
		"data must be a slice or array, got %s",
		e.Kind,
	)
}

// EmptyError is returned when the data given is empty.
type EmptyError struct{}

// Error returns the error message
func (e EmptyError) Error() string {
	return "data must contain at least one element"
}

// InvalidTypeError is returned when the type of the data is not a struct.
type InvalidTypeError struct {
	Kind reflect.Kind
}

// Error returns the error message
func (e InvalidTypeError) Error() string {
	return fmt.Sprintf(
		"data elements must be structs, got %s",
		e.Kind,
	)
}
