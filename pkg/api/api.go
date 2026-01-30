package api

import (
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/ribice/gorsk/pkg/api/auth"
	"github.com/ribice/gorsk/pkg/api/password"
	"github.com/ribice/gorsk/pkg/api/user"
	"github.com/ribice/gorsk/pkg/utl/config"
	"github.com/ribice/gorsk/pkg/utl/jwt"
	"github.com/ribice/gorsk/pkg/utl/middleware"
	"github.com/ribice/gorsk/pkg/utl/query"
	"github.com/ribice/gorsk/pkg/utl/rbac"
	"github.com/ribice/gorsk/pkg/utl/server"
	"github.com/ribice/gorsk/pkg/utl/secure"
	"github.com/ribice/gorsk/pkg/utl/zlog"
)

// New creates new api handlers
func New(e *echo.Echo, cfg *config.Config, log *zlog.Logger,
	ud user.UDB, aut auth.Auth, jwtr jwt.JWT, sec secure.Secure,
	pwd password.PDB, rbac rbac.RBAC) {

	jwtMiddleware := middleware.JWT(cfg.JWT.Secret)
	token := jwtr.NewToken()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Public routes
	e.GET("/status", func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})
	e.GET("/health", server.Health) // Add new health endpoint

	auth.NewHTTP(cfg, log, aut, jwtr, sec, e.Group("/login"), e.Group("/refresh"))

	// Restricted group
	r := e.Group("/v1")
	r.Use(jwtMiddleware)

	user.NewHTTP(log, ud, sec, token, rbac, r.Group("/users"))
	password.NewHTTP(log, pwd, ud, sec, r.Group("/password"))

	// Admin group
	authGroup := r.Group("")
	authGroup.Use(middleware.RBAc(rbac))
	authGroup.Use(middleware.CheckAuth())
	// authGroup.GET("/users/:id", userHandler.View)

	// Private group for swagger
	p := e.Group("")
	p.Use(jwtMiddleware)
	p.Static("/swaggerui", "assets/swaggerui")

	// Assign query parser to all requests
	e.Use(query.Parse)

	// Dummy for jwt expiration. Remove in production
	// e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
	// 	return func(c echo.Context) error {
	// 		auth := c.Request().Header.Get("Authorization")
	// 		if len(auth) > 7 {
	// 			auth = auth[7:]
	// 			tok, _ := jwt.Parse(auth, func(token *jwt.Token) (interface{}, error) {
	// 				return []byte(cfg.JWT.Secret), nil
	// 			})
	// 			if tok != nil {
	// 				claims := tok.Claims.(jwt.MapClaims)
	// 				exp := int64(claims["exp"].(float64))
	// 				log.Info(time.Unix(exp, 0).Sub(time.Now()))
	// 			}
	// 		}
	// 		return next(c)
	// 	}
	// })
}
