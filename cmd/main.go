package main

import (
	"backendService/internals/setup/config"
	"backendService/internals/setup/database"
	"backendService/internals/setup/routes"
	"fmt"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
)

func main() {
	// Loading Configs
	config.LoadConfig()

	// Setup Database
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println("Initializing Database")
		if err := database.InitializeDataBase(config.Config.Database.Type); err != nil {
			panic(err)
		}
	}()

	// Wait for database initialization to complete
	wg.Wait()
	// Setting up Gin
	gin.SetMode(config.Config.App.GinMode)
	app := gin.New()
	app.Use(gin.Recovery())

	// Setting up routes
	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println("Setting up routes")
		routes.SetupAllRoutes(app)

	}()
	wg.Wait()
	// Running the application

	address := ":" + strconv.Itoa(config.Config.App.Port)
	app.Run(address)
}
