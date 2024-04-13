package main

import (
	"backendService/internals/common/logger"
	"backendService/internals/setup/app"
)

// main is the entry point of the application.
func main() {
	logger.Info("main", "main", "Starting Application...")
	app.Start()
}
