package main

import (
	setup "backendService/internals/setup/app"
	"strconv"

	"github.com/gin-gonic/gin"
)

func main() {
	//Loading Configs
	setup.LoadConfig()

	gin.SetMode(setup.Config.App.GinMode)
	app := gin.New()
	app.Use(gin.Recovery())

	//Setup Database
	// database.InitializeDataBase(setup.AppConfig.Database.Type)

	//Setting up routes
	setup.SetupAllRoutes(app)

	//Running the application
	address := ":" + strconv.Itoa(setup.Config.App.Port)
	app.Run(address)
}
