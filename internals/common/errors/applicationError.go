package errors

import (
	"net/http"
)

// ApplicationError represents an application-specific error.
type ApplicationError struct {
	error
	ErrorCode      string      `json:"errorCode"`  // ErrorCode represents the error code associated with the application error.
	Message        string      `json:"message"`    // Message represents the error message associated with the application error.
	HttpStatusCode int         `json:"statusCode"` // HttpStatusCode represents the HTTP status code associated with the application error.
	Err            interface{} `json:"error"`      // Err represents the underlying error associated with the application error.
}

// NewApplicationError creates an instance of ApplicationError with an errorCode, a message, and optional parameters.
// Parameters can specify an HTTP status code and an error object. Defaults to 500 if not provided or incorrect.
func NewApplicationError(errorCode string, message string, parameters ...interface{}) *ApplicationError {
	switch len(parameters) {
	case 0:
		return &ApplicationError{
			ErrorCode:      errorCode,
			Message:        message,
			HttpStatusCode: http.StatusInternalServerError,
		}
	case 1:
		statusCode, ok := parameters[0].(int)
		if !ok {
			statusCode = http.StatusInternalServerError
		}
		return &ApplicationError{
			ErrorCode:      errorCode,
			Message:        message,
			HttpStatusCode: statusCode,
		}
	case 2:
		statusCode, ok := parameters[0].(int)
		if !ok {
			statusCode = http.StatusInternalServerError
		}

		return &ApplicationError{
			ErrorCode:      errorCode,
			Message:        message,
			HttpStatusCode: statusCode,
			Err:            parameters[1],
		}
	default:
		return &ApplicationError{
			ErrorCode:      errorCode,
			Message:        message,
			HttpStatusCode: http.StatusInternalServerError,
		}
	}
}
