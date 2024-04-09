package errors

import (
	"net/http"
)

// UnprocessableEntityError represents an unprocessable entity error.
type UnprocessableEntityError struct {
	ApplicationError
}

// NewUnprocessableEntityError creates an instance of UnprocessableEntityError with an error code, a message, and optional parameters.
// Parameters can specify an HTTP status code and an error object. Defaults to 422 if not provided or incorrect.

func NewUnprocessableEntityError(errorCode string, errors interface{}, parameters ...interface{}) *ApplicationError {
	//Default status code to 422
	statusCode := http.StatusUnprocessableEntity

	if len(parameters) == 0 {
		parameters = append(parameters, statusCode)
	}
	parameters = append(parameters, errors)
	message := "Unprocessable Entity"
	appErr := NewApplicationError(errorCode, message, parameters...)
	return appErr
}
