package errors

import "net/http"

// UnprocessableEntityError represents an unprocessable entity error.
type UnprocessableEntityError struct {
	ApplicationError
}

// NewUnprocessableEntityError creates an instance of UnprocessableEntityError with an error code, a message, and optional parameters.
// Parameters can specify an HTTP status code and an error object. Defaults to 422 if not provided or incorrect.

func NewUnprocessableEntityError(errorCode string, message string, parameters ...interface{}) *UnprocessableEntityError {
	//Default status code to 422
	if parameters[0] == nil {
		parameters[0] = http.StatusUnprocessableEntity
	}

	appErr := NewApplicationError(errorCode, message, parameters...)
	return &UnprocessableEntityError{ApplicationError: *appErr}
}
