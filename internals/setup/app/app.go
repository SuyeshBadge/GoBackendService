package app

import (
	"backendService/internals/common/cache"
	"backendService/internals/common/logger"
	"backendService/internals/setup/config"
	"backendService/internals/setup/database"
	"backendService/internals/setup/server"
)

func Start() {
	// Loading Configs
	config.LoadConfig()

	// Setup Database
	database.InitializeDataBase(config.Config.Database.Type)
	// Setup Cache
	cache.InitializeCacheService()

	// Setting up server
	server := server.NewServer(&config.Config, database.Db)

	// Setting up routes
	SetupAllRoutes(server.App)

	// Running the application
	logger.Info("app", "Start", "Application Started on port", config.Config.App.Port)
	server.Start()

}
