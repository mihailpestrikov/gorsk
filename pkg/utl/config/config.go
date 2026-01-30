package config

import (
	"github.com/kelseyhightower/envconfig"
)

// Config holds configuration data
type Config struct {
	Server   Server
	Postgres Postgres
	JWT      JWT
	App      App
}

// App holds application configuration
type App struct {
	MinPasswordLength int `json:"min_password_length" envconfig:"MIN_PASSWORD_LENGTH" default:"8"`
}

// Server holds server configuration
type Server struct {
	Port         string `json:"port" envconfig:"SERVER_PORT" default:":8080"`
	Debug        bool   `json:"debug" envconfig:"SERVER_DEBUG" default:"false"`
	ReadTimeout  int    `json:"read_timeout" envconfig:"SERVER_READ_TIMEOUT" default:"5"`
	WriteTimeout int    `json:"write_timeout" envconfig:"SERVER_WRITE_TIMEOUT" default:"5"`
	IdleTimeout  int    `json:"idle_timeout" envconfig:"SERVER_IDLE_TIMEOUT" default:"5"`
}

// Postgres holds postgres configuration
type Postgres struct {
	Host     string `json:"host" envconfig:"POSTGRES_HOST" default:"localhost"`
	Port     string `json:"port" envconfig:"POSTGRES_PORT" default:"5432"`
	User     string `json:"user" envconfig:"POSTGRES_USER" default:"gorsk"`
	Password string `json:"password" envconfig:"POSTGRES_PASSWORD" default:"gorsk"`
	DBName   string `json:"db_name" envconfig:"POSTGRES_DBNAME" default:"gorsk"`
	SSLMode  string `json:"ssl_mode" envconfig:"POSTGRES_SSLMODE" default:"disable"`
}

// JWT holds jwt configuration
type JWT struct {
	Secret           string `json:"secret" envconfig:"JWT_SECRET" default:"some_secret"`
	SigningAlgorithm string `json:"signing_algorithm" envconfig:"JWT_SIGNING_ALGORITHM" default:"HS256"`
	Duration         int    `json:"duration" envconfig:"JWT_DURATION" default:"10"`
}

// Load configurations from environment
func Load() (*Config, error) {
	var cfg Config

	err := envconfig.Process("", &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
