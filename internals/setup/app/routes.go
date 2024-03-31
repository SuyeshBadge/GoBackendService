package app

import (
	"backendService/internals/modules/userModule"

	"github.com/gin-gonic/gin"
)

func SetupAllRoutes(app *gin.Engine) {

	userModule.Initialize()
	userModule.UserRouter.SetupRoutes(app)

}
