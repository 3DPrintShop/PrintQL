package errors

import (
	"bytes"
	"fmt"

	pkgerrors "github.com/pkg/errors"
)

// Errors ia a explicit type of a slice of errors.
type Errors []error

// Err returns a single error that represents all of the errors within the collection
func (e Errors) Err() error {
	if len(e) == 0 {
		return nil
	}

	return e
}

// Returns a string that represents all of the errors within the collection
func (e Errors) Error() string {
	var buf bytes.Buffer

	if n := len(e); n == 1 {
		buf.WriteString("1 error: ")
	} else {
		fmt.Fprintf(&buf, "%d errors: ", n)
	}

	for i, err := range e {
		if i != 0 {
			buf.WriteString("; ")
		}

		buf.WriteString(err.Error())
	}

	return buf.String()
}

// Slice returns a slice of errors.
func (e Errors) Slice() []error {
	return []error(e)
}

// New is a convenience method for creating new error types
func New(message string) error {
	return pkgerrors.New(message)
}

// Errorf wraps pkgerrors Errorf
func Errorf(format string, args ...interface{}) error {
	return pkgerrors.Errorf(format, args...)
}
