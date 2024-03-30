package main

import (
	"backendService/internals/setup/config"
	"backendService/internals/setup/database"
	"backendService/internals/setup/routes"
	"strconv"

	"github.com/gin-gonic/gin"
)

// main is the entry point of the application.
// It loads the configurations, initializes the database,
// sets up the Gin framework, sets up the routes,
// and runs the application on the specified port.
func main() {
	// Loading Configs
	config.LoadConfig()

	// Setup Database
	database.InitializeDataBase(config.Config.Database.Type)

	// Setting up Gin
	gin.SetMode(config.Config.App.GinMode)
	app := gin.Default()
	// app.Use(gin.Recovery())

	// Setting up routes
	routes.SetupAllRoutes(app)

	// Running the application
	address := ":" + strconv.Itoa(config.Config.App.Port)
	app.Run(address)
}
