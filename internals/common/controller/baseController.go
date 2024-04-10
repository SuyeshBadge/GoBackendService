package controllers

import (
	"backendService/internals/common/errors"
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type BaseController struct{}

// TransformAndValidate method transforms and validates the given data using the provided DTO struct.
func (c *BaseController) TransformAndValidate(ctx *gin.Context, dtoStruct interface{}) (interface{}, any) {
	if !c.shouldValidate(dtoStruct) {
		return dtoStruct, nil
	}

	if err := ctx.ShouldBindJSON(dtoStruct); err != nil {
		log.Println(err)
		validationErrors := c.extractValidationErrors(err)
		if len(validationErrors) > 0 {
			return nil, errors.NewUnprocessableEntityError("invalid_body", c.newValidationError(validationErrors))
		}
		return nil, err
	}
	validate := validator.New()
	if err := validate.Struct(dtoStruct); err != nil {
		log.Println(err)
		validationErrors := c.extractValidationErrors(err)
		if len(validationErrors) > 0 {
			return nil, errors.NewUnprocessableEntityError("invalid_body", c.newValidationError(validationErrors))
		}
	}

	return dtoStruct, nil
}

// shouldValidate checks if the provided DTO struct needs to be validated.
func (c *BaseController) shouldValidate(dtoStruct interface{}) bool {
	dtoType := reflect.TypeOf(dtoStruct)
	if dtoType == nil {
		return false
	}

	switch dtoType.Kind() {
	case reflect.String, reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Float32, reflect.Float64, reflect.Array, reflect.Slice, reflect.Map:
		return false
	default:
		return true
	}
}

// ValidationErrorData represents a validation error data.
type ValidationErrorData struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ErrValidation is a custom error type for validation errors.
type ErrValidation []ValidationErrorData

func (e ErrValidation) Error() []ValidationErrorData {
	var errorMessages []ValidationErrorData
	for _, validationError := range e {
		errorMessages = append(errorMessages, validationError)

	}
	return errorMessages

}

func (c *BaseController) newValidationError(validationErrors []ValidationErrorData) []ValidationErrorData {
	if len(validationErrors) > 0 {
		return ErrValidation(validationErrors)
	}
	return nil
}

// extractValidationErrors extracts validation errors from the given error and returns a slice of ValidationErrorData.
func (c *BaseController) extractValidationErrors(err error) []ValidationErrorData {
	var validationErrors []ValidationErrorData
	errors, ok := err.(validator.ValidationErrors)
	if ok {
		for _, e := range errors {
			field := e.Field()
			message := c.formatValidationErrorMessage(e, field)
			validationErrors = append(validationErrors, ValidationErrorData{
				Field:   field,
				Message: message,
			})
		}
	} else {
		if t, ok := err.(*json.UnmarshalTypeError); ok {
			validationErrors = append(validationErrors, ValidationErrorData{
				Field:   t.Field,
				Message: c.formatValidationErrorMessage(err, t.Field),
			})
		}

	}
	return validationErrors
}

// formatValidationErrorMessage formats the validation error message based on the tag and field name.
func (c *BaseController) formatValidationErrorMessage(err error, field string) string {
	// Handle unmarshal type errors
	if _, ok := err.(*json.UnmarshalTypeError); ok {

		return fmt.Sprintf("Invalid type for field '%s'", field)
	}

	// Handle other validator errors
	validatorErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return fmt.Sprintf("Invalid value for field '%s'", field)
	}

	for _, validatorError := range validatorErrors {
		if validatorError.Field() == field {
			switch validatorError.Tag() {
			case "required":
				return fmt.Sprintf("%s is required", strings.Title(field))
			case "min":
				return fmt.Sprintf("%s must be at least %s characters long", strings.Title(field), c.getTagValue(validatorError.Tag()))
			case "max":
				return fmt.Sprintf("%s must not exceed %s characters", strings.Title(field), c.getTagValue(validatorError.Tag()))
			case "gte":
				return fmt.Sprintf("%s must be greater than or equal to %s", strings.Title(field), c.getTagValue(validatorError.Tag()))
			case "email":
				return fmt.Sprintf("%s must be a valid email address", strings.Title(field))
			default:
				return fmt.Sprintf("%s is invalid", strings.Title(field))
			}
		}
	}

	return fmt.Sprintf("Invalid value for field '%s'", field)
}

// getTagValue extracts the value from the validation tag.
func (c *BaseController) getTagValue(tag string) string {
	parts := strings.Split(tag, "=")
	if len(parts) > 1 {
		return parts[1]
	}
	return ""
}
