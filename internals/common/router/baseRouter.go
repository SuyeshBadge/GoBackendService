package router

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// BaseRouter is a struct that encapsulates common functionality for handling HTTP requests
type BaseRouter struct {
	Name   string
	Engine *gin.Engine
	group  *gin.RouterGroup
}

// HandlerFunc represents a handler function for processing HTTP requests
type HandlerFunc func(*gin.Context) (interface{}, error)

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

	// Convert middlewares to gin.HandlerFunc
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
			formatSuccessResponse(c, http.StatusOK, data)
		}
	}
}

// formatErrorResponse formats and sends an error response
func formatErrorResponse(c *gin.Context, statusCode int, err any) {
	c.AbortWithStatusJSON(statusCode, gin.H{
		"success":   false,
		"error":     err,
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

// formatSuccessResponse formats and sends a success response
func formatSuccessResponse(c *gin.Context, statusCode int, data interface{}) {
	c.JSON(statusCode, gin.H{
		"success":   true,
		"data":      data,
		"timestamp": time.Now().Format(time.RFC3339),
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
