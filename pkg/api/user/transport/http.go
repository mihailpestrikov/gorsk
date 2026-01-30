package transport

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go-rest-api/pkg/api/user"
	"go-rest-api/pkg/utl/config"
	"go-rest-api/pkg/utl/middleware"
)

type HTTP struct {
	svc user.Service
	mw  *middleware.Middleware
	cfg *config.Config // Add config field
}

// NewHTTP creates new user http service
func NewHTTP(svc user.Service, mw *middleware.Middleware, r *gin.RouterGroup, cfg *config.Config) {
	h := HTTP{svc, mw, cfg} // Initialize cfg

	r.POST("/users", h.create)
	r.GET("/users", h.mw.AuthMiddleware(), h.mw.RbacMiddleware(user.ListRoles), h.list)
	r.GET("/users/:id", h.mw.AuthMiddleware(), h.mw.RbacMiddleware(user.ViewRoles), h.view)
	r.PUT("/users/:id", h.mw.AuthMiddleware(), h.mw.RbacMiddleware(user.EditRoles), h.update)
	r.DELETE("/users/:id", h.mw.AuthMiddleware(), h.mw.RbacMiddleware(user.DeleteRoles), h.delete)
}

type createUserReq struct {
	Username    string `json:"username" binding:"required,min=4"`
	Email       string `json:"email" binding:"required,email"`	
	Password    string `json:"password" binding:"required"` // Removed min=8 binding
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	PhoneNumber string `json:"phone_number"`
	Address     string `json:"address"`
}

func (h *HTTP) create(c *gin.Context) {
	req := new(createUserReq)

	if err := c.BindJSON(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Manual password length validation using the configured value
	if len(req.Password) < h.cfg.MinPasswordLength {
		c.JSON(http.StatusBadRequest, gin.H{"error": "password is too short"})
		return
	}

	res, err := h.svc.Create(c, user.User{Username: req.Username, Email: req.Email, Password: req.Password, FirstName: req.FirstName, LastName: req.LastName, PhoneNumber: req.PhoneNumber, Address: req.Address})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, res)
}

type updateReq struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	PhoneNumber string `json:"phone_number"`
	Address     string `json:"address"`
}

func (h *HTTP) update(c *gin.Context) {
	req := new(updateReq)

	if err := c.BindJSON(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.svc.Update(c, c.Param("id"), user.User{FirstName: req.FirstName, LastName: req.LastName, PhoneNumber: req.PhoneNumber, Address: req.Address})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *HTTP) list(c *gin.Context) {
	// p := c.MustGet("pagination").(*utl.Pagination)

	res, err := h.svc.List(c)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *HTTP) view(c *gin.Context) {
	res, err := h.svc.View(c, c.Param("id"))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *HTTP) delete(c *gin.Context) {
	err := h.svc.Delete(c, c.Param("id"))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}