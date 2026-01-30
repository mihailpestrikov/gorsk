package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"

	"github.com/example/project/pkg/utl/config"
	"github.com/example/project/pkg/utl/jwt"
	"github.com/example/project/pkg/utl/middleware"
	"github.com/example/project/pkg/utl/rbac"
	"github.com/example/project/pkg/utl/server"
	"github.com/example/project/pkg/utl/secure"
	"github.com/example/project/pkg/utl/zlog"
	// ... other imports
)

// New returns the API handler
func New(cfg *config.Config, log *zlog.Logger, db *sqlx.DB) *mux.Router {
	router := mux.NewRouter()

	// Existing middlewares (assuming they are set up similarly with gorilla/mux)
	// If middleware functions wrap http.Handler, they would be applied to subrouters or handlers.
	// For example: router.Use(middleware.Logger(log))
	// Health check endpoint is public and does not require these middlewares.

	// Health check endpoint (public, no authentication required)
	router.HandleFunc("/health", server.Health).Methods(http.MethodGet)

	// ... other existing routes (e.g., /login, /users, etc.)
	// Example of how other routes might be registered with subrouters:
	// authRouter := router.PathPrefix("/auth").Subrouter()
	// auth.New(authRouter, cfg, log, db, jwt.New(cfg.JWT.Secret))

	return router
}
