package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/ribice/gorsk/pkg/utl/server" // Corrected import path
)

// API represents the API structure.
// Assuming this struct already exists or is implicitly handled in main.go
type API struct {
	// Add necessary dependencies here, e.g., Logger, services
	// log *zap.SugaredLogger
}

// Routes sets up the API routes for the application.
// This function needs to be called from the main application to register routes.
// The syntax errors 'expected type, found '{” and 'expected declaration, found '}'
// indicate that routing code was likely placed directly in the package scope,
// outside a function or method.
func (a *API) Routes(r chi.Router) {
	// Global middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	// r.Use(middleware.Logger) // Commented out as project might have custom logging
	r.Use(middleware.Recoverer)

	// Public routes (no authentication required)
	r.Group(func(r chi.Router) {
		r.Get("/health", server.Health) // Register the health check endpoint
	})

	// Example: Protected routes (requires authentication)
	// r.Group(func(r chi.Router) {
	// 	r.Use(authMiddleware.New(jwtService).Auth) // Assuming an auth middleware
	// 	// Add other protected routes here
	// })
}
