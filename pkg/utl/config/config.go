package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// Config holds the configuration for the application
type Config struct {
	Server Server
	DB     DB
	JWT    JWT
	App    App
}

// Server holds server related configurations
type Server struct {
	Port            string `yaml:"port"`
	Timeout         int    `yaml:"timeout"`
	ReadTimeout     int    `yaml:"read_timeout"`
	WriteTimeout    int    `yaml:"write_timeout"`
	IdleTimeout     int    `yaml:"idle_timeout"`
	Debug           bool   `yaml:"debug"`
	MaxHeaderBytes  int    `yaml:"max_header_bytes"`
	SwaggerEnabled  bool   `yaml:"swagger_enabled"`
	PprofEnabled    bool   `yaml:"pprof_enabled"`
	Cors            bool   `yaml:"cors"`
	AllowedOrigins  string `yaml:"allowed_origins"`
	TrustedProxies  string `yaml:"trusted_proxies"`
	CSRF            bool   `yaml:"csrf"`
	CSRFKey         string `yaml:"csrf_key"`
	CSRFHTTPOnly    bool   `yaml:"csrf_http_only"`
	CSRFSecure      bool   `yaml:"csrf_secure"`
	CSRFSameSite    string `yaml:"csrf_same_site"`
	CSRFPath        string `yaml:"csrf_path"`
	CSRFMaxAge      int    `yaml:"csrf_max_age"`
	SessionCookie   bool   `yaml:"session_cookie"`
	SessionName     string `yaml:"session_name"`
	SessionKey      string `yaml:"session_key"`
	SessionHTTPOnly bool   `yaml:"session_http_only"`
	SessionSecure   bool   `yaml:"session_secure"`
	SessionSameSite string `yaml:"session_same_site"`
	SessionPath     string `yaml:"session_path"`
	SessionMaxAge   int    `yaml:"session_max_age"`
}

// DB holds database related configurations
type DB struct {
	Dialect         string `yaml:"dialect"`
	URL             string `yaml:"url"`
	Schema          string `yaml:"schema"`
	MaxOpenConns    int    `yaml:"max_open_conns"`
	MaxIdleConns    int    `yaml:"max_idle_conns"`
	ConnMaxLifetime int    `yaml:"conn_max_lifetime"`
	LogQueries      bool   `yaml:"log_queries"`
}

// JWT holds JWT related configurations
type JWT struct {
	Secret        string `yaml:"secret"`
	Duration      int    `yaml:"duration"`
	RefreshSecret string `yaml:"refresh_secret"`
	RefreshDuration int  `yaml:"refresh_duration"`
}

// App holds application specific configurations
type App struct {
	AppURL  string  `yaml:"app_url"`
	Version Version `yaml:"version"` // Added Version field
}

// Version holds version and name information
type Version struct {
	Name    string `yaml:"name" json:"name"`       // Added json tag
	Version string `yaml:"version" json:"version"` // Added json tag
}

// New returns a new Config instance
func New(file string) (*Config, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg);
		if err != nil {
			return nil, err
		}
	return &cfg, nil
}
