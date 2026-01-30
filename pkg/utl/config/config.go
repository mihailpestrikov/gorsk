package config

type Config struct {
	// Server settings
	ServerPort         string `envconfig:"server_port" default:":8080"`
	Debug              bool   `envconfig:"debug" default:"false"`
	ReadTimeoutSeconds int    `envconfig:"read_timeout_seconds" default:"5"`
	WriteTimeoutSeconds int   `envconfig:"write_timeout_seconds" default:"10"`
	MinPasswordLength  int    `envconfig:"min_password_length" default:"8"`

	// Database settings
	DBHost     string `envconfig:"db_host" default:"localhost"`
	DBPort     string `envconfig:"db_port" default:"5432"`
	DBUser     string `envconfig:"db_user" default:"postgres"`
	DBPassword string `envconfig:"db_password" default:"postgres"`
	DBName     string `envconfig:"db_name" default:"go_rest_api"`
	DBSchema   string `envconfig:"db_schema" default:"public"`
	SSLMode    string `envconfig:"ssl_mode" default:"disable"`

	// JWT settings
	JWTSecretKey           string `envconfig:"jwt_secret_key" default:"secret"`
	JWTAccessTokenLifeSpan int    `envconfig:"jwt_access_token_life_span" default:"15"`  // in minutes
	JWTRefreshTokenLifeSpan int   `envconfig:"jwt_refresh_token_life_span" default:"1440"` // in minutes
}