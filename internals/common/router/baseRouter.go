package router

import (
	"backendService/internals/common/errors"
	"backendService/internals/setup/server"
	"log"
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
				// log.Println("inside handleWrapper2", err)
				formatErrorResponse(c, http.StatusInternalServerError, err)
			}
		}()
		data, err := handler(c)
		if err != nil {
			// log.Println("inside handleWrapper", err)
			formatErrorResponse(c, http.StatusInternalServerError, err)
		} else {
			formatSuccessResponse(c, http.StatusOK, data.Data, data.Message)
		}
	}
}

// formatErrorResponse formats and sends an error response
func formatErrorResponse(c *gin.Context, statusCode int, err interface{}) {

	var errorCode, message string

	if appErr, ok := err.(*errors.ApplicationError); ok {
		statusCode = appErr.HttpStatusCode
		errorCode = appErr.ErrorCode
		message = appErr.Message
	} else {
		goErr := goError.Wrap(err, 2)
		stackTrace := goErr.ErrorStack()
		log.Println(stackTrace)
		errorCode = "internal_server_error"
		if server.Server.Config.App.Env == "development" {
			message = err.(error).Error()
		} else {
			message = "Something went wrong. Please try again later."
		}
		err = errors.NewApplicationError(errorCode, message, http.StatusInternalServerError)

	}

	c.JSON(statusCode, gin.H{
		"error":     gin.H{"errorCode": errorCode, "message": message},
		"success":   false,
		"timestamp": time.Now().Format(time.RFC3339),
	})
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
