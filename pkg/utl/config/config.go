package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/spf13/viper"
)

// Config holds all application configuration
type Config struct {
	Port              string        `yaml:"port"`
	Env               string        `yaml:"env"`
	LogFile           string        `yaml:"log_file"`
	LogHTTP           bool          `yaml:"log_http"`
	DB                Database      `yaml:"database"`
	JWT               JWT           `yaml:"jwt"`
	ReadTimeout       time.Duration `yaml:"read_timeout"`
	WriteTimeout      time.Duration `yaml:"write_timeout"`
	GraceTimeout      time.Duration `yaml:"grace_timeout"`
	MinPasswordLength int           `yaml:"min_password_length"` // New field
}

// Database holds database configuration
type Database struct {
	Type     string `yaml:"type"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
	SSLMode  string `yaml:"ssl_mode"`
}

// JWT holds jwt configuration
type JWT struct {
	Secret    string        `yaml:"secret"`
	Duration  time.Duration `yaml:"duration"`
	Refresh   time.Duration `yaml:"refresh"`
	MaxRefresh time.Duration `yaml:"max_refresh"`
}

// LoadConfig loads configuration from file or environment
func LoadConfig(path string) (*Config, error) {
	v.AddConfigPath(path)
	v.SetConfigName("conf")
	v.SetConfigType("yaml")

	v.SetDefault("port", "8080")
	v.SetDefault("env", "development")
	v.SetDefault("log_file", "app.log")
	v.SetDefault("log_http", true)
	v.SetDefault("jwt.duration", "1h")
	v.SetDefault("jwt.refresh", "24h")
	v.SetDefault("jwt.max_refresh", "720h")
	v.SetDefault("read_timeout", "5s")
	v.SetDefault("write_timeout", "10s")
	v.SetDefault("grace_timeout", "5s")
	v.SetDefault("min_password_length", 8) // Set default for new field

	v.AutomaticEnv() // Read from environment variables

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println("Config file not found, using defaults and environment variables.")
		} else {
			return nil, err
		}
	}

	cfg := &Config{}
	if err := v.Unmarshal(cfg); err != nil {
		return nil, err
	}

	// Environment variable overrides for MinPasswordLength (optional, but good practice)
	if minLenStr := os.Getenv("MIN_PASSWORD_LENGTH"); minLenStr != "" {
		if minLen, err := strconv.Atoi(minLenStr); err == nil {
			cfg.MinPasswordLength = minLen
		} else {
			log.Printf("Warning: Invalid MIN_PASSWORD_LENGTH environment variable: %v, using default/config value.", err)
		}
	}

	return cfg, nil
}
