package errors

import (
	"backendService/internals/common/logger"
)

type InternalServerError struct {
	ApplicationError
}

func NewInternalServerError(message string, cause interface{}) *ApplicationError {
	//log the error
	logger.Error("errors", "NewInternalServerError", "Internal Server Error", message)
	message = "Something went wrong internally. Please try again later."

	return NewApplicationError("internal_server_error", message)
}
