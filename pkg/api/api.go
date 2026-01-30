package api

import (
	"net/http"

	"github.com/go-pg/pg/v9"
	"github.com/labstack/echo/v4"
	apiAuth "github.com/ribice/gorsk/pkg/api/auth"
	apiUser "github.com/ribice/gorsk/pkg/api/user"
	"github.com/ribice/gorsk/pkg/utl/config"
	"github.com/ribice/gorsk/pkg/utl/jwt"
	"github.com/ribice/gorsk/pkg/utl/middleware"
	"github.com/ribice/gorsk/pkg/utl/rbac"
	"github.com/ribice/gorsk/pkg/utl/secure"
	"github.com/ribice/gorsk/pkg/utl/server"
	"github.com/rs/zerolog"
)

// API groups all application services
type API struct {
	Auth apiAuth.Service
	User apiUser.Service
}

// New creates new api
func New(cfg *config.Config, db *pg.DB, rbac RBAC, sec Secure, jwt JWT, log *zerolog.Logger) (*echo.Echo, error) {
	srv := server.New(cfg, log)
	authMiddleware := middleware.JWT(jwt.Secret())

	// Public routes
	srv.Echo.GET("/health", srv.Health())
	srv.Echo.POST("/login", apiAuth.Login(cfg, jwt, srv, db, sec, log))
	srv.Echo.GET("/refresh/:token", apiAuth.Refresh(cfg, jwt, srv, db, log))

	// Private routes
	v1 := srv.Echo.Group("/v1", authMiddleware)

	// Auth
	v1.PATCH("/password/:id", apiAuth.ChangePassword(cfg, srv, db, sec, log))

	// Users
	v1.POST("/users", apiUser.Create(cfg, srv, db, rbac, sec, log))
	v1.GET("/users", apiUser.List(cfg, srv, db, rbac, log))
	v1.GET("/users/:id", apiUser.View(cfg, srv, db, rbac, log))
	v1.PATCH("/users/:id", apiUser.Update(cfg, srv, db, rbac, log))
	v1.DELETE("/users/:id", apiUser.Delete(cfg, srv, db, rbac, log))

	//Swagger
	srv.Echo.Static("/swaggerui", "./assets/swaggerui")
	srv.Echo.GET("/swagger.json", func(c echo.Context) error {
		return c.File("assets/swaggerui/swagger.json")
	})

	return srv.Echo, nil
}

type (
	// JWT represents JWT interface
	JWT interface {
		GenerateToken(apiAuth.User) (string, error)
		ParseToken(string) (*jwt.AuthClaims, error)
		Secret() string
	}

	// RBAC represents RBAC interface
	RBAC interface {
		UserHasPermission(jwt.AuthClaims, string, string) error
		UserHasRole(jwt.AuthClaims, string) error
		EnforceUser(jwt.AuthClaims, string, ...string) error
		Enforce(jwt.AuthClaims, string, string, ...string) error
		EnforceErr(jwt.AuthClaims, string, string, ...string) error
	}

	// Secure represents security interface
	Secure interface {
		Hash(string) (string, error)
		CompareHash(string, string) error
	}
)
