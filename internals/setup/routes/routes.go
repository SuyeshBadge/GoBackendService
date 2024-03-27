package routes

import (
	userModule "backendService/internals/modules/userModule/routes"

	"github.com/gin-gonic/gin"
)

func SetupAllRoutes(app *gin.Engine) {
	userModule.UserRoutes.SetupRoutes(app)
}
