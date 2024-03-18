package main

import (
	setup "backendService/internals/setup/app"
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	app := gin.New()
	app.Use(gin.Recovery())

	setup.LoadConfig()

	fmt.Println("Starting the server...")
	fmt.Print("Environment: ", setup.AppConfig.App.GinMode, "\n")

	setup.SetupAllRoutes(app)

	app.Run(":8100")
}
