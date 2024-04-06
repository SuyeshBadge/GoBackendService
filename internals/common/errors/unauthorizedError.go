package errors

import "net/http"

type UnauthorizedError struct {
	ApplicationError
}

// NewUnauthorizedError creates an instance of UnauthorizedError with an error code, a message, and optional parameters.
// Parameters can specify an HTTP status code and an error object. Defaults to 401 if not provided or incorrect.

func NewUnauthorizedError(errorCode string, message string, parameters ...interface{}) *UnauthorizedError {
	//Default status code to 401
	if parameters[0] == nil {
		parameters[0] = http.StatusUnauthorized
	}

	appErr := NewApplicationError(errorCode, message, parameters...)
	return &UnauthorizedError{ApplicationError: *appErr}
}
