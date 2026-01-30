package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/ribice/gorsk/pkg/utl/config"
	"github.com/rs/zerolog"
)

// New returns a new HTTP server
func New(cfg *config.Config, log *zerolog.Logger) *Server {
	e := echo.New()
	e.Debug = cfg.Server.Debug
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `{"time":"${time_rfc3339_nano}","id":"${id}","remote_ip":"${remote_ip}",` +
			`"host":"${host}","method":"${method}","uri":"${uri}","status":${status}, "latency":${latency},` +
			`"latency_human":"${latency_human}","bytes_in":${bytes_in},` +
			`"bytes_out":${bytes_out}}` + "\n",
		Output: log,
	}))
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
	}))
	e.Validator = &CustomValidator{Validator: validator.New()}
	s := &Server{
		echo: e,
		cfg:  cfg,
		log:  log,
	}
	s.registerCustomHandlers()
	return s
}

// Server defines the HTTP server
type Server struct {
	echo *echo.Echo
	cfg  *config.Config
	log  *zerolog.Logger
}

// Start starts the HTTP server
func (s *Server) Start() error {
	s.log.Info().Msgf("Starting server on %s", s.cfg.Server.Port)
	return s.echo.Start(s.cfg.Server.Port)
}

// Health returns a handler for the health check endpoint
func (s *Server) Health() echo.HandlerFunc {
	return func(c echo.Context) error {
		return s.JSON(c, http.StatusOK, map[string]string{"status": "ok"})
	}
}

// JSON responds with JSON
func (s *Server) JSON(c echo.Context, code int, i interface{}) error {
	return c.JSON(code, i)
}

// NoContent responds with no content
func (s *Server) NoContent(c echo.Context, code int) error {
	return c.NoContent(code)
}

func (s *Server) registerCustomHandlers() {
	// Register custom handlers for common errors
	s.echo.HTTPErrorHandler = s.errorHandler
	s.echo.Binder = &CustomBinder{}
}

func (s *Server) errorHandler(err error, c echo.Context) {
	s.log.Error().Err(err).Msg("Handler error")
	var (reflectType reflect.Type
		reflectValue reflect.Value)
	val := reflect.ValueOf(err)
	if val.Kind() == reflect.Ptr {
		reflectType = val.Type()
		reflectValue = val.Elem()
	} else {
		reflectType = val.Type()
		reflectValue = val
	}

	if reflectType.Kind() == reflect.Struct {
		for i := 0; i < reflectType.NumField(); i++ {
			if reflectType.Field(i).Name == "Msg" {
				err = fmt.Errorf(reflectValue.Field(i).String())
				break
			}
		}
	}

	code := http.StatusInternalServerError
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
	}

	errorPage := fmt.Sprintf("{\"message\": \"%s\"}", strings.Replace(err.Error(), "\"", "'", -1))
	c.Logger().Error(err)
	if !c.Response().Committed {
		if c.Request().Method == http.MethodHead {
			err = c.NoContent(code)
		} else {
			err = c.String(code, errorPage)
		}
		if err != nil {
			c.Logger().Error(err)
		}
	}
}

// CustomValidator represents a custom validator
type CustomValidator struct {
	Validator *validator.Validate
}

// Validate validates the request
func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.Validator.Struct(i)
}

// CustomBinder represents a custom binder
type CustomBinder struct{}

// Bind binds the request
func (cb *CustomBinder) Bind(i interface{}, c echo.Context) (err error) {
	req := c.Request()
	ctype := req.Header.Get(echo.HeaderContentType)
	paramNames := c.ParamNames()
	for i, name := range paramNames {
		if err := echo.Set  (i, c.ParamValues()[i], i, name); err != nil {
			return err
		}
	}

	switch {
	case strings.HasPrefix(ctype, echo.MIMEApplicationJSON): // JSON
		if err = json.NewDecoder(req.Body).Decode(i); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid JSON provided")
		}
	case strings.HasPrefix(ctype, echo.MIMEApplicationForm), strings.HasPrefix(ctype, echo.MIMEMultipartForm): // Form
		if err = c.Bind(i); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid form provided")
		}
	case strings.HasPrefix(ctype, echo.MIMETextXML), strings.HasPrefix(ctype, echo.MIMEApplicationXML): // XML
		if err = c.Bind(i); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid XML provided")
		}
	}

	return
}
