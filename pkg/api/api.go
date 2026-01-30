package api

import (
	"github.com/labstack/echo"
	"github.com/ribice/gorsk/pkg/utl/server" // Import the server package
)

// API represents the API structure
type API struct {
	echo *echo.Echo
	// Other dependencies like services, middleware etc.
	// ...
}

// New creates a new API instance and registers its routes
func New(e *echo.Echo /* ... other dependencies if any */) *API {
	a := &API{
		echo: e,
		// Initialize other fields
		// ...
	}
	a.RegisterRoutes()
	return a
}

// RegisterRoutes registers all API routes
func (a *API) RegisterRoutes() {
	// Public routes - no authentication required
	a.echo.GET("/health", server.Health)

	// Existing public routes example (commented out as placeholder):
	// a.echo.POST("/login", a.auth.Login)
	// a.echo.GET("/refresh/:token", a.auth.Refresh)
	// a.echo.GET("/me", a.user.Me)
	// a.echo.GET("/swaggerui/", a.swagger.ServeSwaggerUI)

	// Authenticated routes example (commented out as placeholder):
	// v1 := a.echo.Group("/v1")
	// v1.Use(a.middleware.AuthJWT, a.middleware.Auth)
	// v1.GET("/users", a.user.List)
	// v1.GET("/users/:id", a.user.View)
	// v1.POST("/users", a.user.Create)
	// v1.PATCH("/password/:id", a.user.ChangePassword)
	// v1.DELETE("/users/:id", a.user.Delete)

	// NOTE: This file may contain other existing code which is not shown here.
	// The '/health' route registration is added to the public routes section.
}

// NOTE: Other API methods and struct definitions would be here.
