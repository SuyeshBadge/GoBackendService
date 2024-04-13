package router

import (
	"backendService/internals/common/errors"
	"backendService/internals/common/logger"
	"backendService/internals/setup/server"
	"net/http"
	"time"

	goError "github.com/go-errors/errors"

	"github.com/gin-gonic/gin"
)

// BaseRouter is a struct that encapsulates common functionality for handling HTTP requests
type BaseRouter struct {
	Name   string
	Engine *gin.Engine
	group  *gin.RouterGroup
}

type Response struct {
	Data    any    `json:"data"`
	Message string `json:"message"`
}

// HandlerFunc represents a handler function for processing HTTP requests
type HandlerFunc func(*gin.Context) (Response, interface{})

// NewBaseRouter initializes and returns a new instance of the BaseRouter struct
func NewBaseRouter(name string, engine *gin.Engine) *BaseRouter {
	return &BaseRouter{
		Name:   name,
		Engine: engine,
		group:  engine.Group("/"),
	}
}

// Handle registers routes with the specified HTTP method, path, and handler function(s)
func (br *BaseRouter) Handle(method, path string, handlers ...HandlerFunc) {
	middlewares, lastHandler := getMiddlewaresAndLastHandler(handlers)
	wrappedLastHandler := handleWrapper(lastHandler)
	ginMiddlewares := make([]gin.HandlerFunc, len(middlewares))
	for i, mw := range middlewares {
		ginMiddlewares[i] = handleWrapper(mw)
	}
	br.group.Handle(method, path, append(ginMiddlewares, wrappedLastHandler)...)
}

// getMiddlewaresAndLastHandler extracts middlewares and the final handler function
func getMiddlewaresAndLastHandler(handlers []HandlerFunc) ([]HandlerFunc, HandlerFunc) {
	lastIndex := len(handlers) - 1
	return handlers[:lastIndex], handlers[lastIndex]
}

// handleWrapper wraps the final handler function with Gin middleware for error handling and response formatting
func handleWrapper(handler HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				formatErrorResponse(c, http.StatusInternalServerError, err)
			}
		}()
		data, err := handler(c)
		if err != nil {
			formatErrorResponse(c, http.StatusInternalServerError, err)
		} else {
			formatSuccessResponse(c, http.StatusOK, data.Data, data.Message)
		}
	}
}

// formatErrorResponse formats and sends an error response
func formatErrorResponse(c *gin.Context, statusCode int, err interface{}) {

	var errorCode, message string

	var validationError interface{}

	// Check if the error is an ApplicationError
	if appErr, ok := err.(*errors.ApplicationError); ok {
		// If so, set the status code, error code, message, and validation error from the ApplicationError
		statusCode = appErr.HttpStatusCode
		errorCode = appErr.ErrorCode
		message = appErr.Message
		validationError = appErr.Err
	} else {
		// If not, wrap the error with goError and get the stack trace
		goErr := goError.Wrap(err, 2)
		stackTrace := goErr.ErrorStack()
		// Log the stack trace
		logger.Error("router", "baseRouter", "formatErrorResponse", err, stackTrace)
		// Set a generic error code and message
		errorCode = "internal_server_error"
		// If in development environment, show the actual error message
		if server.Server.Config.App.Env == "development" {
			message = err.(error).Error()
		} else {
			// Otherwise, show a generic error message to the user
			message = "Something went wrong. Please try again later."
		}
		// Create a new ApplicationError with the generic error code and message
		err = errors.NewApplicationError(errorCode, message, http.StatusInternalServerError)

	}

	// Check if the status code is 422 Unprocessable Entity
	if statusCode == http.StatusUnprocessableEntity {
		// If so, return a JSON response with detailed error information including validation errors
		c.JSON(statusCode, gin.H{
			"error": gin.H{
				"errorCode": errorCode,       // Error code for the specific error
				"message":   message,         // Error message describing the issue
				"errors":    validationError, // Detailed validation error messages
			},
			"success":   false,                           // Indicate the operation was not successful
			"timestamp": time.Now().Format(time.RFC3339), // Timestamp of the error occurrence
		})
	} else {
		// For other status codes, return a JSON response with general error information
		c.JSON(statusCode, gin.H{
			"error":     gin.H{"errorCode": errorCode, "message": message}, // Error code and message
			"success":   false,                                             // Indicate the operation was not successful
			"timestamp": time.Now().Format(time.RFC3339),                   // Timestamp of the error occurrence
		})
	}

}

// formatSuccessResponse formats and sends a success response
func formatSuccessResponse(c *gin.Context, statusCode int, data interface{}, message string) {
	c.JSON(statusCode, gin.H{
		"success":   true,
		"data":      data,
		"timestamp": time.Now().Format(time.RFC3339),
		"message":   message,
	})
}

// Group creates a new router group relative to the current router's path and applies provided handlers to it
func (br *BaseRouter) Group(relPath string, handlers ...gin.HandlerFunc) *BaseRouter {
	newGroup := br.group.Group(relPath, handlers...)
	return &BaseRouter{
		Name:   br.Name,
		Engine: br.Engine,
		group:  newGroup,
	}
}

// GET registers a route with the GET HTTP method
func (br *BaseRouter) GET(path string, handlers ...HandlerFunc) {
	br.Handle(http.MethodGet, path, handlers...)
}

// POST registers a route with the POST HTTP method
func (br *BaseRouter) POST(path string, handlers ...HandlerFunc) {
	br.Handle(http.MethodPost, path, handlers...)
}

// PUT registers a route with the PUT HTTP method
func (br *BaseRouter) PUT(path string, handlers ...HandlerFunc) {
	br.Handle(http.MethodPut, path, handlers...)
}

// DELETE registers a route with the DELETE HTTP method
func (br *BaseRouter) DELETE(path string, handlers ...HandlerFunc) {
	br.Handle(http.MethodDelete, path, handlers...)
}

// PATCH registers a route with the PATCH HTTP method
func (br *BaseRouter) PATCH(path string, handlers ...HandlerFunc) {
	br.Handle(http.MethodPatch, path, handlers...)
}

// OPTIONS registers a route with the OPTIONS HTTP method
func (br *BaseRouter) OPTIONS(path string, handlers ...HandlerFunc) {
	br.Handle(http.MethodOptions, path, handlers...)
}
