package api

import (
	"encoding/json" // Added import for JSON encoding
	"net/http"

	"github.com/go-chi/chi"
	"gopkg.in/go-playground/validator.v9"
	"gorsk.io/gorsk/pkg/api/auth"
	"gorsk.io/gorsk/pkg/api/password"
	"gorsk.io/gorsk/pkg/api/user"
	"gorsk.io/gorsk/pkg/utl/config" // Added import for config
	"gorsk.io/gorsk/pkg/utl/jwt"
	"gorsk.io/gorsk/pkg/utl/middleware"
	"gorsk.io/gorsk/pkg/utl/rbac"
	"gorsk.io/gorsk/pkg/utl/server" // Assumed for server.Respond
)

// API provides the application resources
type API struct {
	config    *config.Config
	router    *chi.Mux
	// Services
	auth      auth.Service
	psw       password.Service
	usr       user.Service
	jwt       jwt.Service
	rbac      rbac.Service
	mw        middleware.Service
	validator *validator.Validate
}

// New creates a new instance of API
func New(cfg *config.Config,
	auth auth.Service,
	psw password.Service,
	usr user.Service,
	jwt jwt.Service,
	rbac rbac.Service,
	mw middleware.Service,
	v *validator.Validate,
) *API {
	return &API{
		config:    cfg,
		router:    chi.NewRouter(),
		auth:      auth,
		psw:       psw,
		usr:       usr,
		jwt:       jwt,
		rbac:      rbac,
		mw:        mw,
		validator: v,
	}
}

// Router returns the API router
func (a *API) Router() *chi.Mux {
	a.router.Use(a.mw.Logger)
	// Public routes
	a.registerPublicRoutes(a.router)

	// Private routes
	a.router.Group(func(r chi.Router) {
		r.Use(a.mw.Auth, a.mw.RBAC)
		a.registerPrivateRoutes(r)
	})

	return a.router
}

// registerPublicRoutes adds the public routes
func (a *API) registerPublicRoutes(r *chi.Mux) {
	r.Mount("/auth", a.auth.Router())
	r.Mount("/password", a.psw.Router())
	// Add the version endpoint here
	r.Get("/version", a.version)
}

// registerPrivateRoutes adds the private routes
func (a *API) registerPrivateRoutes(r *chi.Mux) {
	r.Mount("/users", a.usr.Router())
}

// version returns the API version and name
// @Summary Get API version
// @Description Get current API version and build information
// @Tags Public
// @Accept json
// @Produce json
// @Success 200 {object} config.Version "API version information"
// @Router /version [get]
func (a *API) version(w http.ResponseWriter, r *http.Request) {
	server.Respond(w, http.StatusOK, a.config.App.Version)
}
