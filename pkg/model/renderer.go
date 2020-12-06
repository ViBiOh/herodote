package model

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

var (
	// ErrInvalid occurs when checks fails
	ErrInvalid = errors.New("invalid")

	// ErrNotFound occurs when somehint is not found
	ErrNotFound = errors.New("not found")

	// ErrInternalError occurs when shit happens
	ErrInternalError = errors.New("internal error")
)

// Message for render
type Message struct {
	Level   string
	Content string
}

// ParseMessage parses messages from request
func ParseMessage(r *http.Request) Message {
	return Message{
		Level:   r.URL.Query().Get("messageLevel"),
		Content: r.URL.Query().Get("messageContent"),
	}
}

// NewSuccessMessage create a success message
func NewSuccessMessage(content string) Message {
	return Message{
		Level:   "success",
		Content: content,
	}
}

// NewErrorMessage create a error message
func NewErrorMessage(content string) Message {
	return Message{
		Level:   "error",
		Content: content,
	}
}

// ConcatError concat errors to a single string
func ConcatError(errs []error) error {
	if len(errs) == 0 {
		return nil
	}

	values := make([]string, len(errs))
	for index, err := range errs {
		values[index] = err.Error()
	}

	return errors.New(strings.Join(values, ", "))
}

// WrapInvalid wraps given error with invalid err
func WrapInvalid(err error) error {
	return fmt.Errorf("%s: %w", err, ErrInvalid)
}

// WrapInternal wraps given error with internal err
func WrapInternal(err error) error {
	return fmt.Errorf("%s: %w", err, ErrInternalError)
}

// WrapNotFound wraps given error with not found err
func WrapNotFound(err error) error {
	return fmt.Errorf("%s: %w", err, ErrNotFound)
}
