package errors

import "log"

type InternalServerError struct {
	ApplicationError
}

func NewInternalServerError(message string, cause interface{}) *ApplicationError {
	//log the error
	log.Fatalln(message+" :: ", cause)

	message = "Something went wrong internally. Please try again later."

	return NewApplicationError("internal_server_error", message)
}
