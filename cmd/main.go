package main

import (
	setup "backendService/internals/setup/app"

	"github.com/gin-gonic/gin"
)

func main() {
	app := gin.New()
	app.Use(gin.Recovery())

	setup.SetupAllRoutes(app)

	app.Run(":8100")
}
