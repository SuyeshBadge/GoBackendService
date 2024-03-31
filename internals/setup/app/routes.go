package app

import (
	"backendService/internals/modules/userModule"
	"backendService/internals/setup/database"

	"github.com/gin-gonic/gin"
)

func SetupAllRoutes(app *gin.Engine) {
	app.GET("/", func(c *gin.Context) {
		databaseName := database.Db.Name()
		c.JSON(200, gin.H{
			"message": databaseName,
		})
	})

	userModule.Initialize()
	userModule.UserRouter.SetupRoutes(app)

}
