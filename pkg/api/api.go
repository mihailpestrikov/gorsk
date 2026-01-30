package api

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
	"github.com/ribice/gorsk/pkg/api/auth"
	"github.com/ribice/gorsk/pkg/api/password"
	"github.com/ribice/gorsk/pkg/api/user"
	"github.com/ribice/gorsk/pkg/utl/config"
	"github.com/ribice/gorsk/pkg/utl/jwt"
	"github.com/ribice/gorsk/pkg/utl/middleware"
	"github.com/ribice/gorsk/pkg/utl/server"
	"github.com/ribice/gorsk/pkg/utl/zlog"
)

// New returns a new API
func New(cfg *config.Config, jwtSvc jwt.Service, log zlog.Logger,
	auth auth.Service, mw *middleware.Middleware,
	user user.Service, pass password.Service) *echo.Echo {

	e := echo.New()
	e.Use(mw.CORS)
	e.Use(mw.Logger())

	// Public routes
	// login
	e.POST("/login", auth.Login)
	// refresh token
	e.GET("/refresh/:token", auth.Refresh)
	// health check
	e.GET("/health", server.Health)

	// Secured routes
	g := e.Group("/v1")
	g.Use(mw.Auth)

	// users
	g.GET("/users", user.List)
	g.POST("/users", user.Create)
	g.GET("/users/:id", user.View)
	g.PUT("/users/:id", user.Update)
	g.DELETE("/users/:id", user.Delete)
	g.PATCH("/password/:id", pass.Change)

	// Swagger
	e.Static("/swaggerui", cfg.Server.Swagger.DocsPath)
	e.GET("/swagger.json", func(c echo.Context) error {
		return c.File(fmt.Sprintf("%s/swagger.json", cfg.Server.Swagger.DocsPath))
	})

	return e
}
