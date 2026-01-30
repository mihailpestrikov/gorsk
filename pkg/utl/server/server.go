package server

import (
	"net/http"
	"time"

	"github.com/rs/cors"
	"github.com/urfave/negroni"

	"app/pkg/utl/config" // Assuming this path based on project structure
	"app/pkg/utl/middleware" // Import the new middleware package
	"app/pkg/utl/zlog"   // Assuming this path based on project structure
)

// Server holds common server functionalities and dependencies.
type Server struct {
	Config *config.Config
	Logger *zlog.Logger
}

// New creates a new server instance with the given configuration and logger.
func New(cfg *config.Config, logger *zlog.Logger) *Server {
	return &Server{
		Config: cfg,
		Logger: logger,
	}
}

// Start sets up the negroni middleware chain and starts the HTTP server.
// It expects a configured http.Handler (typically a router like Mux or Chi).
func (s *Server) Start(router http.Handler) error {
	n := negroni.New()

	// Standard middlewares
	n.Use(negroni.NewRecovery()) // Catches panics and returns a 500 error
	n.Use(middleware.RequestID()) // Add Request ID middleware early in the chain
	n.Use(negroni.NewLogger())   // Logs requests; should be after RequestID to include it in logs

	// CORS middleware configuration
	n.Use(cors.New(cors.Options{
		AllowedOrigins: []string{"*"}, // Adjust as per your security requirements
		AllowedMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{
			"Authorization",
			"Content-Type",
			middleware.HeaderXRequestID, // Allow X-Request-ID from client
		},
		ExposedHeaders: []string{
			middleware.HeaderXRequestID, // Expose X-Request-ID to clients
		},
		Debug: s.Config.Server.Debug, // Assuming Config.Server.Debug exists
	}))

	// The main application router is the final handler
	n.UseHandler(router)

	srv := &http.Server{
		Addr:         s.Config.Server.Port, // Assuming Config.Server.Port exists
		Handler:      n,
		ReadTimeout:  time.Duration(s.Config.Server.ReadTimeout) * time.Second,  // Assuming Config.Server.ReadTimeout exists
		WriteTimeout: time.Duration(s.Config.Server.WriteTimeout) * time.Second, // Assuming Config.Server.WriteTimeout exists
		IdleTimeout:  time.Duration(s.Config.Server.IdleTimeout) * time.Second,  // Assuming Config.Server.IdleTimeout exists
	}

	s.Logger.Info().Msgf("Server listening on port %s", s.Config.Server.Port)
	return srv.ListenAndServe()
}
