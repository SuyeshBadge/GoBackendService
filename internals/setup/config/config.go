package config

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

// DBConfig holds the database configuration values
type Database struct {
	Type     string `mapstructure:"type"`
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Database string `mapstructure:"database"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

// AppConfig holds the application configuration values
type ApplicationConfig struct {
	Port     int    `mapstructure:"port"`
	Env      string `mapstructure:"env"`
	LogLevel string `mapstructure:"log_level"`
	GinMode  string `mapstructure:"gin_mode"`
}

// Config holds the overall configuration
type AppConfig struct {
	Database Database          `mapstructure:"database"`
	App      ApplicationConfig `mapstructure:"app"`
}

var Config AppConfig

// LoadConfig loads the configuration from environment variables and config files
func LoadConfig() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	//set config file name
	env := os.Getenv("APP_ENV")

	fmt.Println("Loading config for environment: ", env)

	// Set the config file name
	configDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting current directory: %v", err)
	}
	configDir = configDir + "/configs"
	viper.SetConfigName("default")
	viper.SetConfigType("json")
	viper.AddConfigPath(configDir)

	// Read the config file
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	if env == "development" {
		viper.SetConfigName("development")
	} else if env == "production" {
		viper.SetConfigName("production")
	} else {
		viper.SetConfigName("default")
	}

	if err := viper.MergeInConfig(); err != nil {
		log.Fatalf("Error merging %s config file: %v", env, err)
	}
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// Unmarshal config values into struct
	if err := viper.Unmarshal(&Config); err != nil {
		log.Fatalf("Unable to decode config into struct: %v", err)
	}
}
