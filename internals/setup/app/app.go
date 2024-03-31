package app

import (
	"backendService/internals/setup/config"
	"backendService/internals/setup/database"
	"backendService/internals/setup/server"
)

func Start() {
	// Loading Configs
	config.LoadConfig()

	// Setup Database
	database.InitializeDataBase(config.Config.Database.Type)

	// Setting up server
	server := server.NewServer(&config.Config, database.Db)

	// Setting up routes
	SetupAllRoutes(server.App)

	// Running the application
	server.Start()

}
