package router

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type BaseRouter struct {
	RouterName string
	Router     *gin.Engine
	group      *gin.RouterGroup
}

func NewBaseRouter(routerName string, router *gin.Engine) *BaseRouter {
	return &BaseRouter{
		RouterName: routerName,
		Router:     router,
	}
}

type HandlerFunc func(c *gin.Context) (interface{}, error)

// HandleFunc registers a route with the given method, path and handler function.
func (br *BaseRouter) Handle(method string, path string, handler ...interface{}) {
	middlewares, lastHandler := getMiddlewaresAndLastHandler(handler...)

	br.Router.Handle(method, path, append(middlewares, handleWrapper(lastHandler))...)
}

// get all middlewares and last handler
func getMiddlewaresAndLastHandler(handler ...interface{}) ([]gin.HandlerFunc, HandlerFunc) {
	//get all middlewares
	middlewares := []gin.HandlerFunc{}
	for _, h := range handler {
		if m, ok := h.(gin.HandlerFunc); ok {
			middlewares = append(middlewares, m)
		}
	}

	//get last handler
	lastHandler := handler[len(handler)-1].(HandlerFunc)

	return middlewares, lastHandler
}

func handleWrapper(handler HandlerFunc) func(c *gin.Context) {

	return func(c *gin.Context) {
		//call handler without switch case

		data, err := handler(c)
		if err != nil {
			formatErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
		formatSuccessResponse(c, http.StatusOK, data)

	}
}

func formatErrorResponse(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, gin.H{
		"success":   false,
		"error":     message,
		"timestamp": time.Now().Unix(),
	})
}

func formatSuccessResponse(c *gin.Context, statusCode int, data interface{}) {
	c.JSON(statusCode, gin.H{
		"success":   true,
		"data":      data,
		"timestamp": time.Now().Unix(),
	})
}

func (br *BaseRouter) Group(relativePath string, handlers ...gin.HandlerFunc) *BaseRouter {
	br.group = br.Router.Group(relativePath, handlers...)
	return br

}

func (br *BaseRouter) GET(relativePath string, handlers ...interface{}) {
	middlewares, lastHandler := getMiddlewaresAndLastHandler(handlers...)

	temp := append(middlewares, handleWrapper(lastHandler))
	log.Println("temp", temp)
	//if group is not nil, then add handlers to group else add to router
	if br.group != nil {
		br.group.GET(relativePath, append(middlewares, handleWrapper(lastHandler))...)
	} else {
		br.Router.GET(relativePath, append(middlewares, handleWrapper(lastHandler))...)
	}
}
