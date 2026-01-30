package api

import (
	"my_project/pkg/utl/server" // Assuming the correct import path based on go.mod
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	// Other specific API package imports, e.g., auth, user, password, company, role, etc.
)

// API defines the API application service.
// It encapsulates the application's dependencies and routes.
type API struct {
	server *server.Server // Dependency for server-level utilities like health checks
	// authService    auth.Service    // Example: Authentication service
	// userService    user.Service    // Example: User management service
	// passwordService password.Service // Example: Password operations
	// ... other service dependencies
}

// New creates a new instance of the API service.
// It takes all necessary services and dependencies as arguments.
func New(srv *server.Server /* authSvc auth.Service, userSvc user.Service, passwordSvc password.Service */) *API {
	return &API{
		server: srv,
		// authService:    authSvc,
		// userService:    userSvc,
		// passwordService: passwordSvc,
	}
}

// Router configures and returns the main HTTP router for the API.
func (a *API) Router() *chi.Mux {
	r := chi.NewRouter()

	// Global Middlewares
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	// r.Use(middleware.URLFormat) // Example: remove trailing slashes

	// Public Routes - no authentication required
	r.Group(func(r chi.Router) {
		r.Get("/health", a.server.Health) // Add the health check endpoint
		// r.Post("/login", a.authService.LoginHandler) // Example: existing public route
		// r.Post("/register", a.userService.RegisterHandler) // Example: existing public route
	})

	// Authenticated Routes - require authentication middleware
	// r.Group(func(r chi.Router) {
	// 	// Assuming `a.authService` provides middleware for JWT verification and authentication
	// 	// r.Use(a.authService.VerifyJWT)
	// 	// r.Use(a.authService.AuthenticateUser)
	//
	// 	r.Route("/users", func(r chi.Router) {
	// 		// r.Get("/", a.userService.ListUsersHandler)
	// 		// r.Post("/", a.userService.CreateUserHandler)
	// 		// r.Get("/{id}", a.userService.GetUserHandler)
	// 		// r.Put("/{id}", a.userService.UpdateUserHandler)
	// 		// r.Delete("/{id}", a.userService.DeleteUserHandler)
	// 	})
	//
	// 	// Add other protected routes here
	// })

	return r
}
