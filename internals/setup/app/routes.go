package app

import (
	"backendService/internals/modules/authModule"
	"backendService/internals/modules/userModule"

	"github.com/gin-gonic/gin"
)

func SetupAllRoutes(app *gin.Engine) {

	userModule.Initialize()
	authModule.Initialize()

	userModule.UserRouter.SetupRoutes(app)
	authModule.AuthRouter.SetupRoutes(app)

}
