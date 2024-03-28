package routes

import (
	userModule "backendService/internals/modules/userModule/routes"
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
	userModule.UserRoutes.SetupRoutes(app)
}
