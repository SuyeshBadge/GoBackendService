package config

import (
	"fmt"
	"os"
	"strings"

	"backendService/internals/common/logger"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

var Config AppConfig

// LoadConfig loads the configuration from environment variables and config files
func LoadConfig() {
	// Load environment variables from .env file
	if err := godotenv.Load(".env"); err != nil {
		logger.Fatal("config", "LoadConfig", "loadEnv", err)
	}

	// Set config file name
	env := os.Getenv("APP_ENV")

	logger.Info("config", "LoadConfig", "setConfigFile", fmt.Sprintf("Using config file: %s.json", env))

	// Set the config file directory
	configDir, err := os.Getwd()
	if err != nil {
		logger.Fatal("config", "LoadConfig", "getWorkingDir", err)
	}
	configDir += "/configs"

	viper.SetConfigType("json")
	viper.AddConfigPath(configDir)

	// Set the base config file
	viper.SetConfigName("default")
	if err := viper.ReadInConfig(); err != nil {
		logger.Fatal("config", "LoadConfig", "readConfigFile", err)
	}

	// Set the environment-specific config file
	configName := "default"
	if env == "development" || env == "production" {
		configName = env
	}
	viper.SetConfigName(configName)
	if err := viper.MergeInConfig(); err != nil {
		logger.Fatal("config", "LoadConfig", "mergeConfigFile", err)
	}

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// Unmarshal config values into struct
	if err := viper.Unmarshal(&Config); err != nil {
		logger.Fatal("config", "LoadConfig", "unmarshal", err)
	}
}
