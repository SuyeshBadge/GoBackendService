package config

// Database holds the database configuration values
type Database struct {
	Type     string `mapstructure:"type"`
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Database string `mapstructure:"database"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

// ApplicationConfig holds the application configuration values
type ApplicationConfig struct {
	Port     int    `mapstructure:"port"`
	Env      string `mapstructure:"env"`
	LogLevel string `mapstructure:"log_level"`
	GinMode  string `mapstructure:"gin_mode"`
}

type CacheConfig struct {
	Hostname string `mapstructure:"hostname"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	Database int    `mapstructure:"database"`
}

// AppConfig holds the overall configuration
type AppConfig struct {
	Database Database          `mapstructure:"database"`
	App      ApplicationConfig `mapstructure:"app"`
	Cache    CacheConfig       `mapstructure:"cache"`
}
