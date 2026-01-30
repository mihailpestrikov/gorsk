package transport

import (
	"encoding/json"
	"fmt" // For error formatting
	"net/http"

	"your_project/pkg/api/user"          // Assuming user service interface and errors are here
	"your_project/pkg/utl/config"        // New import for config
	"your_project/pkg/utl/server"        // For common HTTP server utilities
	"github.com/go-chi/chi"
	"github.com/go-playground/validator/v10"
)

// UserService defines methods the handler expects to be implemented by a user service.
type UserService interface {
	Create(u *user.User) (*user.User, error)
	// ... other user-related methods
}

// HTTP represents the user HTTP handler.
type HTTP struct {
	svc       UserService
	cfg       *config.Config       // Added config field
	validator *validator.Validate
}

// NewHTTP creates a new user HTTP handler.
func NewHTTP(svc UserService, cfg *config.Config, r *chi.Mux) { // Added cfg parameter
	h := HTTP{
		svc:       svc,
		cfg:       cfg,              // Assign config
		validator: validator.New(),
	}
	r.Post("/users", h.create)
	// Add other user routes here, e.g.,
	// r.Get("/users/{id}", h.view)
	// r.Put("/users/{id}", h.update)
	// r.Delete("/users/{id}", h.delete)
}

// createUserReq defines the structure for a user creation request.
type createUserReq struct {
	FirstName string `json:"first_name" validate:"required,max=50"`
	LastName  string `json:"last_name" validate:"max=50"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required"` // Removed min=8 tag
	Phone     string `json:"phone" validate:"max=20"`
	Address   string `json:"address" validate:"max=100"`
	Active    bool   `json:"active"`
	RoleID    int    `json:"role_id" validate:"required"`
	CompanyID int    `json:"company_id" validate:"required"`
}

// create handles user creation requests.
func (h *HTTP) create(w http.ResponseWriter, r *http.Request) {
	req := new(createUserReq)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		server.RespondError(w, http.StatusBadRequest, err)
		return
	}

	// Manual password length validation using the configurable value
	if len(req.Password) < h.cfg.MinPasswordLength {
		server.RespondError(w, http.StatusBadRequest,
			fmt.Errorf("password must be at least %d characters long", h.cfg.MinPasswordLength))
		return
	}

	if err := h.validator.Struct(req); err != nil {
		server.RespondError(w, http.StatusBadRequest, err)
		return
	}

	// Create a user object from the request
	usr := &user.User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Password:  req.Password, // Password will be hashed by the service layer
		Phone:     req.Phone,
		Address:   req.Address,
		Active:    req.Active,
		RoleID:    req.RoleID,
		CompanyID: req.CompanyID,
	}

	createdUser, err := h.svc.Create(usr)
	if err != nil {
		server.RespondError(w, http.StatusInternalServerError, err) // Or specific error codes for service errors
		return
	}

	server.RespondJSON(w, http.StatusCreated, createdUser)
}
