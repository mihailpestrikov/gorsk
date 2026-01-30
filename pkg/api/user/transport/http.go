package transport

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ribice/gorsk/pkg/api/user/service"
	"github.com/ribice/gorsk/pkg/utl/config"
	"github.com/ribice/gorsk/pkg/utl/middleware"
	"github.com/ribice/gorsk/pkg/utl/model"
	"github.com/ribice/gorsk/pkg/utl/zlog"
)

type ( // Request and response structures
	// Create request
	create struct {
		Username  string        `json:"username" validate:"required,min=4,max=50"`
		Password  string        `json:"password" validate:"required,max=250"`
		FirstName string        `json:"first_name" validate:"required,min=2,max=50"`
		LastName  string        `json:"last_name" validate:"required,min=2,max=50"`
		Address   string        `json:"address" validate:"required,min=2,max=250"`
		Phone     string        `json:"phone" validate:"required,min=2,max=50"`
		Email     string        `json:"email" validate:"required,email"`
		Roles     []model.UserRole `json:"roles" validate:"required,min=1"`
	}

	// Update request
	update struct {
		ID        int           `json:"id" validate:"required,min=1"`
		FirstName string        `json:"first_name" validate:"required,min=2,max=50"`
		LastName  string        `json:"last_name" validate:"required,min=2,max=50"`
		Address   string        `json:"address" validate:"required,min=2,max=250"`
		Phone     string        `json:"phone" validate:"required,min=2,max=50"`
		Email     string        `json:"email" validate:"required,email"`
		Roles     []model.UserRole `json:"roles" validate:"required,min=1"`
	}

	// Change password request
	changePassword struct {
		OldPassword string `json:"old_password" validate:"required"`
		NewPassword string `json:"new_password" validate:"required,max=250"`
	}
)

// H represents user http transport service
type H struct {
	svc    service.Service
	sec    service.Auth
	log    *zlog.Logger
	config *config.Config
}

// NewHTTP creates new user http transport service
func NewHTTP(svc service.Service, sec service.Auth, log *zlog.Logger, cfg *config.Config) *H {
	return &H{svc: svc, sec: sec, log: log, config: cfg}
}

// Create creates a new user account
func (h *H) Create(c echo.Context) error {
	req := new(create)

	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, zlog.Error{Message: err.Error()})
	}

	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, zlog.Error{Message: err.Error()})
	}

	if len(req.Password) < h.config.App.MinPasswordLength {
		return c.JSON(http.StatusBadRequest, zlog.Error{Message: fmt.Sprintf("password must be at least %d characters long", h.config.App.MinPasswordLength)})
	}

	id, err := h.svc.Create(c, model.User{
		Username:  req.Username,
		Password:  req.Password,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Address:   req.Address,
		Phone:     req.Phone,
		Email:     req.Email,
		Roles:     req.Roles,
	})

	if err != nil {
		return c.JSON(http.StatusInternalServerError, zlog.Error{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"id": id})
}

// List returns list of users
func (h *H) List(c echo.Context) error {
	users, err := h.svc.List(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, zlog.Error{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, users)
}

// View returns a user's details
func (h *H) View(c echo.Context) error {
	id, err := h.svc.View(c, c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, zlog.Error{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, id)
}

// Delete deletes a user
func (h *H) Delete(c echo.Context) error {
	if err := h.svc.Delete(c, c.Param("id")); err != nil {
		return c.JSON(http.StatusInternalServerError, zlog.Error{Message: err.Error()})
	}
	return c.NoContent(http.StatusOK)
}

// Update updates user's contact information
func (h *H) Update(c echo.Context) error {
	req := new(update)

	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, zlog.Error{Message: err.Error()})
	}

	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, zlog.Error{Message: err.Error()})
	}

	usr, err := h.svc.Update(c, model.User{
		ID:        req.ID,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Address:   req.Address,
		Phone:     req.Phone,
		Email:     req.Email,
		Roles:     req.Roles,
	})

	if err != nil {
		return c.JSON(http.StatusInternalServerError, zlog.Error{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, usr)
}

// ChangePassword changes user's password
func (h *H) ChangePassword(c echo.Context) error {
	req := new(changePassword)

	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, zlog.Error{Message: err.Error()})
	}

	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, zlog.Error{Message: err.Error()})
	}

	if len(req.NewPassword) < h.config.App.MinPasswordLength {
		return c.JSON(http.StatusBadRequest, zlog.Error{Message: fmt.Sprintf("new password must be at least %d characters long", h.config.App.MinPasswordLength)})
	}

	usr := c.Get(middleware.UserKey).(*model.User)

	if err := h.svc.ChangePassword(c, usr.ID, req.OldPassword, req.NewPassword);
		err != nil {
		return c.JSON(http.StatusInternalServerError, zlog.Error{Message: err.Error()})
	}

	return c.NoContent(http.StatusOK)
}
