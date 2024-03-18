package setup

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

// DBConfig holds the database configuration values
type DBConfig struct {
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

// Database holds the different database configurations
type Database struct {
	Postgres DBConfig `mapstructure:"postgres"`
	MySQL    DBConfig `mapstructure:"mysql"`
	SQLite   DBConfig `mapstructure:"sqlite"`
}

// Config holds the overall configuration
type Config struct {
	Database Database          `mapstructure:"database"`
	App      ApplicationConfig `mapstructure:"app"`
}

var AppConfig Config

// LoadConfig loads the configuration from environment variables and config files
func LoadConfig() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Set the environment variable prefix
	viper.SetEnvPrefix("app")
	viper.AutomaticEnv()

	//set config file name
	env := os.Getenv("APP_ENV")

	fmt.Println("Loading config for environment: ", env)

	if env == "dev" {
		viper.SetConfigName("development")
	} else if env == "prod" {
		viper.SetConfigName("production")
	} else {
		viper.SetConfigName("default")
	}

	// Set the config file name
	configDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting current directory: %v", err)
	}
	configDir = configDir + "/configs"

	viper.AddConfigPath(configDir)

	// Read the config file
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	// Unmarshal config values into struct
	if err := viper.Unmarshal(&AppConfig); err != nil {
		log.Fatalf("Unable to decode config into struct: %v", err)
	}
}
