package errors

import "net/http"

// NotFoundError represents a not found error.
type NotFoundError struct {
	ApplicationError
}

// NewNotFoundError creates an instance of NotFoundError with an error code, a message, and optional parameters.
// Parameters can specify an HTTP status code and an error object. Defaults to 404 if not provided or incorrect.
func NewNotFoundError(errorCode string, message string, parameters ...interface{}) *ApplicationError {
	//Default status code to 404
	if parameters[0] == nil {
		parameters[0] = http.StatusNotFound
	}

	appErr := NewApplicationError(errorCode, message, parameters...)
	return appErr
}
