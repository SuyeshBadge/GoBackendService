package main

import (
	setup "backendService/internals/setup/app"
	"strconv"

	"github.com/gin-gonic/gin"
)

func main() {
	//Loading Configs
	setup.LoadConfig()

	gin.SetMode(setup.AppConfig.App.GinMode)
	app := gin.New()
	app.Use(gin.Recovery())

	//Setting up routes
	setup.SetupAllRoutes(app)

	//Running the application
	address := ":" + strconv.Itoa(setup.AppConfig.App.Port)
	app.Run(address)
}
