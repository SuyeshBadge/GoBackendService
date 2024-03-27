package main

import (
	"backendService/internals/setup/config"
	"backendService/internals/setup/database"
	"backendService/internals/setup/routes"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

func main() {
	//Loading Configs
	config.LoadConfig()

	gin.SetMode(config.Config.App.GinMode)
	app := gin.New()
	app.Use(gin.Recovery())

	//Setup Database
	database.InitializeDataBase(config.Config.Database.Type)

	fmt.Println("Database connected", database.Db.Name())
	//Setting up routes
	routes.SetupAllRoutes(app)

	//Running the application
	address := ":" + strconv.Itoa(config.Config.App.Port)
	app.Run(address)
}
