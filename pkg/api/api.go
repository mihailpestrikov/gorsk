package api

import (
	"github.com/gin-gonic/gin"
	"github.com/ribice/gorsk/pkg/api/auth"
	"github.com/ribice/gorsk/pkg/api/password"
	"github.com/ribice/gorsk/pkg/api/user"
	"github.com/ribice/gorsk/pkg/utl/config"
	"github.com/ribice/gorsk/pkg/utl/jwt"
	"github.com/ribice/gorsk/pkg/utl/middleware"
	"github.com/ribice/gorsk/pkg/utl/server"
)

// Router registers all the routes
func Router(r *gin.Engine, cfg *config.Config, jwtService jwt.Service) {
	// Public endpoints
	r.POST("/login", auth.Login(jwtService))
	r.GET("/health", server.Health) // Add health check endpoint

	// Authenticated endpoints
	authGroup := r.Group("")
	authGroup.Use(middleware.Auth(jwtService))
	{
		authGroup.POST("/password/reset", password.Reset)
		authGroup.GET("/users", user.List)
		authGroup.POST("/users", user.Create)
		authGroup.GET("/users/:id", user.View)
		authGroup.PUT("/users/:id", user.Update)
		authGroup.DELETE("/users/:id", user.Delete)
	}
}
