package api

import (
	"github.com/labstack/echo"
	"github.com/ribice/gorsk/pkg/api/auth"
	"github.com/ribice/gorsk/pkg/api/password"
	"github.com/ribice/gorsk/pkg/api/user"
	"github.com/ribice/gorsk/pkg/utl/config"
	"github.com/ribice/gorsk/pkg/utl/jwt"
	"github.com/ribice/gorsk/pkg/utl/middleware"
	"github.com/ribice/gorsk/pkg/utl/rbac"
	"github.com/ribice/gorsk/pkg/utl/server" // Import the server package for Health handler
	"github.com/ribice/gorsk/pkg/utl/zlog"
)

// Boot creates a new instance of the Echo API and registers all routes.
func Boot(e *echo.Echo, cfg *config.Config, j *jwt.JWT, log *zlog.Logger) error {
	// Initialize RBAC for middleware
	rbacService := rbac.New()

	// Public routes: no authentication required
	// Serve Swagger UI assets
	e.Static("/swaggerui", "assets/swaggerui")
	// Serve the swagger.json file
	e.File("/swagger.json", "assets/swaggerui/swagger.json")

	// Register public authentication related routes (e.g., login, refresh token)
	auth.NewHTTP(auth.NewService(nil, j, cfg), log).RegisterPublic(e)
	// Register public password related routes (e.g., password reset request)
	password.NewHTTP(password.NewService(nil, cfg), j, log).RegisterPublic(e)

	// Add the /health endpoint here as a public route
	e.GET("/health", server.Health)

	// Versioned API group with authentication and authorization middleware
	v1 := e.Group("/v1")
	v1.Use(middleware.JWT(j), middleware.RBAC(rbacService))

	// Register protected v1 routes
	user.NewHTTP(user.NewService(nil, cfg), j, log).RegisterPrivate(v1)
	// Add other protected API services here as needed

	return nil
}
