package controllers

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type BaseController struct{}

// TransformAndValidate method transforms and validates the given data using the provided DTO struct.
func (c *BaseController) TransformAndValidate(ctx *gin.Context, dtoStruct interface{}) (interface{}, error) {
	if !c.shouldValidate(dtoStruct) {
		return dtoStruct, nil
	}

	validate := validator.New()
	err := validate.Struct(dtoStruct)
	if err != nil {
		validationErrors := extractValidationErrors(err)
		// ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
		// 	"errors": validationErrors,
		// })
		return nil, ErrValidation(validationErrors)
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

func (e ErrValidation) Error() string {
	var errorMessages []string
	for _, validationError := range e {
		errorMessages = append(errorMessages, fmt.Sprintf("Field: %s, Message: %s", validationError.Field, validationError.Message))
	}
	return strings.Join(errorMessages, "\n")
}

// extractValidationErrors extracts validation errors from the given error and returns a slice of ValidationErrorData.
func extractValidationErrors(err error) []ValidationErrorData {
	var validationErrors []ValidationErrorData
	errors := err.(validator.ValidationErrors)
	for _, e := range errors {
		message := formatValidationErrorMessage(e.Tag(), e.Field())
		validationErrors = append(validationErrors, ValidationErrorData{
			Field:   e.Field(),
			Message: message,
		})
	}
	return validationErrors
}

// formatValidationErrorMessage formats the validation error message based on the tag and field name.
func formatValidationErrorMessage(tag, field string) string {
	switch tag {
	case "required":
		return fmt.Sprintf("%s is required", strings.Title(field))
	case "min":
		return fmt.Sprintf("%s must be at least %s characters long", strings.Title(field), getTagValue(tag))
	case "max":
		return fmt.Sprintf("%s must not exceed %s characters", strings.Title(field), getTagValue(tag))
	case "gte":
		return fmt.Sprintf("%s must be greater than or equal to %s", strings.Title(field), getTagValue(tag))
	case "email":
		return fmt.Sprintf("%s must be a valid email address", strings.Title(field))
	default:
		return fmt.Sprintf("%s is invalid", strings.Title(field))
	}
}

// getTagValue extracts the value from the validation tag.
func getTagValue(tag string) string {
	parts := strings.Split(tag, "=")
	if len(parts) > 1 {
		return parts[1]
	}
	return ""
}
