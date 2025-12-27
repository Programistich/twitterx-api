package apperror

import (
	"errors"
	"fmt"
	"net/http"
)

// ValidationError represents invalid input from the client
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error: %s - %s", e.Field, e.Message)
}

// NotFoundError represents a resource that was not found
type NotFoundError struct {
	Resource string
	ID       string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("%s '%s' not found", e.Resource, e.ID)
}

// UpstreamError represents an error from an upstream service (Nitter, FxTwitter)
type UpstreamError struct {
	Service    string
	StatusCode int
	Message    string
	Err        error
}

func (e *UpstreamError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s error: %s: %v", e.Service, e.Message, e.Err)
	}
	return fmt.Sprintf("%s error: %s (status: %d)", e.Service, e.Message, e.StatusCode)
}

func (e *UpstreamError) Unwrap() error {
	return e.Err
}

// HTTPStatusCode returns the appropriate HTTP status code for the error
func HTTPStatusCode(err error) int {
	var validationErr *ValidationError
	var notFoundErr *NotFoundError
	var upstreamErr *UpstreamError

	switch {
	case errors.As(err, &validationErr):
		return http.StatusBadRequest
	case errors.As(err, &notFoundErr):
		return http.StatusNotFound
	case errors.As(err, &upstreamErr):
		return http.StatusBadGateway
	default:
		return http.StatusInternalServerError
	}
}
