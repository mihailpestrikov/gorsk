package api

import (
	"github.com/go-chi/chi/v5" // Assuming chi router is used
	"github.com/ribice/gorsk/pkg/api/auth"
	"github.com/ribice/gorsk/pkg/api/password"
	"github.com/ribice/gorsk/pkg/api/user"
	authMW "github.com/ribice/gorsk/pkg/utl/middleware/auth"   // Correct specific auth middleware import
	secureMW "github.com/ribice/gorsk/pkg/utl/middleware/secure" // Correct specific secure middleware import
	"github.com/ribice/gorsk/pkg/utl/server" // Import the server package for the Health handler
)

// Option represents a functional option
type Option func(*API)

// API represents the API structure
type API struct {
	Auth   auth.Service
	Pass   password.Service
	User   user.Service
}

// New creates new api struct
func New(opts ...Option) *API {
	api := &API{}
	for _, opt := range opts {
		opt(api)
	}
	return api
}

// RegisterRoutes registers all application routes
func RegisterRoutes(r chi.Router, api *API) {
	// Add the public health check endpoint
	r.Get("/health", server.Health)

	// Existing base middleware (if any) could be applied here
	r.Use(secureMW.Headers) // Example, assuming secure middleware is used

	// Group for authenticated routes
	r.Group(func(r chi.Router) {
		r.Use(authMW.Authenticate) // Apply authentication middleware
		auth.Routes(r, api.Auth)
		password.Routes(r, api.Pass)
		user.Routes(r, api.User)
	})
}

// WithAuthService is a functional option for Auth service
func WithAuthService(s auth.Service) Option {
	return func(api *API) {
		api.Auth = s
	}
}

// WithPasswordService is a functional option for Password service
func WithAuthService(s password.Service) Option {
	return func(api *API) {
		api.Pass = s
	}
}

// WithUserService is a functional option for User service
func WithUserService(s user.Service) Option {
	return func(api *API) {
		api.User = s
	}
}
