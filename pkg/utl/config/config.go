package config

import (
	"log"
	"os"

	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	App  App  `yaml:"app"`
	DB   DB   `yaml:"db"`
	JWT  JWT  `yaml:"jwt"`
	Mail Mail `yaml:"mail"`
}

type App struct {
	Environment       string `yaml:"environment"`
	MinPasswordLength int    `yaml:"min_password_length"`
}

type DB struct {
	Dialect        string `yaml:"dialect"`
	URL            string `yaml:"url"`
	MigrationURL   string `yaml:"migration_url"`
	MaxOpenConns   int    `yaml:"max_open_conns"`
	MaxIdleConns   int    `yaml:"max_idle_conns"`
	ConnMaxLifetime int    `yaml:"conn_max_lifetime"`
}

type JWT struct {
	Secret        string `yaml:"secret"`
	Duration      int    `yaml:"duration"`
	Refresh       int    `yaml:"refresh"`
	RefreshTime   int    `yaml:"refresh_time"`
	SigningMethod string `yaml:"signing_method"`
}

type Mail struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
	User string `yaml:"user"`	Password string `yaml:"password"`
	From string `yaml:"from"`
}

// LoadConfig loads configuration from file or environment variables
func LoadConfig(configPath string) (*Config, error) {
	config := Config{}

	bytes, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	if err = yaml.Unmarshal(bytes, &config);
		err != nil {
		return nil, err
	}

	env := os.Getenv("ENVIRONMENT_NAME")
	if env == ""
		env = "development"

	log.Printf("Loading \"%s\" environment configuration", env)

	return &config, nil
}
