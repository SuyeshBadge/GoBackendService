package errors

import "net/http"

type BadRequestError struct {
	ApplicationError
}

// NewBadRequestError creates an instance of BadRequestError with an error code, a message, and optional parameters.
// Parameters can specify an HTTP status code and an error object. Defaults to 400 if not provided or incorrect.
func NewBadRequestError(errorCode string, message string, parameters ...interface{}) *BadRequestError {
	//Default status code to 400
	if parameters[0] == nil {
		parameters[0] = http.StatusBadRequest
	}

	appErr := NewApplicationError(errorCode, message, parameters...)
	return &BadRequestError{ApplicationError: *appErr}
}
