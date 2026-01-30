package transport

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/ribice/gorsk/pkg/api/user"
	"github.com/ribice/gorsk/pkg/utl/middleware"
	"github.com/ribice/gorsk/pkg/utl/model"
)

// HTTP represents user http transport
type HTTP struct {
	service user.Service
}

// NewHTTP creates new user http transport
func NewHTTP(svc user.Service, mw *middleware.Middleware) *HTTP {
	return &HTTP{svc}
}

// RegisterHandlers registers handlers
func (h *HTTP) RegisterHandlers(e *echo.Echo, mw *middleware.Middleware) {
	r := e.Group("/v1/users")
	r.Use(mw.JWT)
		r.GET("", h.list)
		r.GET("/:id", h.view)
		r.POST("", h.create, mw.RBAC(model.AdminRole))
		r.PATCH("/:id", h.update, mw.RBAC(model.AdminRole))
		r.DELETE("/:id", h.delete, mw.RBAC(model.AdminRole))

	p := e.Group("/v1/password")
	p.Use(mw.JWT)
		p.PATCH("/:id", h.changePassword)
}

type create struct {
	Username    string `json:"username" validate:"required,min=2,max=50"`
	Password    string `json:"password" validate:"required,min_password,max=60"`
	Email       string `json:"email" validate:"required,email"`
	FirstName   string `json:"first_name" validate:"required,min=2,max=50"`
	LastName    string `json:"last_name" validate:"required,min=2,max=50"`
	PhoneNumber string `json:"phone_number" validate:"required,min=3,max=15"`
	Address     string `json:"address" validate:"required,min=2,max=100"`
	RoleID      int    `json:"role_id" validate:"required,min=1"`
	CompanyID   int    `json:"company_id" validate:"required,min=1"`
	LocationID  int    `json:"location_id" validate:"required,min=1"`
}

// Create user
func (h *HTTP) create(c echo.Context) error {
	req := new(create)

	if err := c.Bind(req); err != nil {
		return err
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	res, err := h.service.Create(c, model.User{Username: req.Username, Password: req.Password, Email: req.Email, FirstName: req.FirstName, LastName: req.LastName, PhoneNumber: req.PhoneNumber, Address: req.Address, RoleID: req.RoleID, CompanyID: req.CompanyID, LocationID: req.LocationID})
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, res)
}

type update struct {
	FirstName   string `json:"first_name,omitempty" validate:"min=2,max=50"`
	LastName    string `json:"last_name,omitempty" validate:"min=2,max=50"`
	PhoneNumber string `json:"phone_number,omitempty" validate:"min=3,max=15"`
	Address     string `json:"address,omitempty" validate:"min=2,max=100"`
	RoleID      int    `json:"role_id,omitempty" validate:"min=1"`
	CompanyID   int    `json:"company_id,omitempty" validate:"min=1"`
	LocationID  int    `json:"location_id,omitempty" validate:"min=1"`
}

// Update user
func (h *HTTP) update(c echo.Context) error {
	req := new(update)

	if err := c.Bind(req); err != nil {
		return err
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err
	}

	if err := h.service.Update(c, req.toModel(id)); err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}

func (u *update) toModel(id int) *model.User {
	return &model.User{
		Base: model.Base{ID: id},
		FirstName:   u.FirstName,
		LastName:    u.LastName,
		PhoneNumber: u.PhoneNumber,
		Address:     u.Address,
		RoleID:      u.RoleID,
		CompanyID:   u.CompanyID,
		LocationID:  u.LocationID,
	}
}

type changePassword struct {
	OldPassword string `json:"old_password" validate:"required,min_password,max=60"`
	NewPassword string `json:"new_password" validate:"required,min_password,max=60"`
}

// ChangePassword changes user's password
func (h *HTTP) changePassword(c echo.Context) error {
	req := new(changePassword)

	if err := c.Bind(req); err != nil {
		return err
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err
	}

	if err := h.service.ChangePassword(c, id, req.OldPassword, req.NewPassword); err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}

// View user
func (h *HTTP) view(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err
	}

	res, err := h.service.View(c, id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, res)
}

// List users
func (h *HTTP) list(c echo.Context) error {
	res, err := h.service.List(c, c.QueryParams())
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, res)
}

// Delete user
func (h *HTTP) delete(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err
	}

	if err := h.service.Delete(c, id); err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}
