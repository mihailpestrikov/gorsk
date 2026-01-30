package server

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/ribice/gorsk/pkg/utl/config"
	"github.com/rs/zerolog"
	gopkg.in/go-playground/validator.v8"
)

// Server contains server configurations
type Server struct {
	*echo.Echo
}

// New returns new HTTP server
func New(cfg *config.Config, logger *zerolog.Logger) *Server {
	e := echo.New()
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Output: logger,
	}))
	e.Use(middleware.Secure())
	e.Use(middleware.CORS())

	v := validator.New()
	v.RegisterValidation("min_password", func(fl validator.FieldLevel) bool {
		return len(fl.Field().String()) >= cfg.App.MinPasswordLength
	})
	e.Validator = &customValidator{validator: v}
	e.HTTPErrorHandler = newHTTPError(logger)

	return &Server{e}
}

type customValidator struct {
	validator *validator.Validate
}

// Validate validates all fields in a struct
func (cv *customValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func newHTTPError(logger *zerolog.Logger) echo.HTTPErrorHandler {
	return func(err error, c echo.Context) {
		reqID := c.Request().Header.Get(echo.HeaderXRequestID)
		var customErr *echo.HTTPError
		if he, ok := err.(*echo.HTTPError); ok {
			customErr = he
		} else if e, ok := err.(validator.ValidationErrors); ok {
			customErr = echo.NewHTTPError(http.StatusBadRequest, e.Error())
		} else {
			customErr = echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		if c.Response().Committed {
			return
		}

		if customErr.Code >= 500 {
			logger.Err(customErr).Str("request_id", reqID).Send()
		} else if customErr.Code >= 400 {
			logger.Warn().Err(customErr).Str("request_id", reqID).Send()
		}

		if c.Request().Method == http.MethodHead {
			err = c.NoContent(customErr.Code)
		} else {
			err = c.JSON(customErr.Code, customErr)
		}
		if err != nil {
			logger.Err(err).Str("request_id", reqID).Msg("http error handler failed")
		}
	}
}
