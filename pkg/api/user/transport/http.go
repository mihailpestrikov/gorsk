package transport

import (
	"net/http"
	"strconv" // For parsing query parameters

	"github.com/gin-gonic/gin"
	"github.com/user/repo/pkg/api/user"       // Import user package for the service interface and domain models
	"github.com/user/repo/pkg/utl/query"      // Import query package for pagination parameters
	"github.com/user/repo/pkg/utl/middleware" // Assuming middleware for handling errors or other aspects
)

// ListResponse contains a list of users and the total count.
type ListResponse struct {
	Users      []user.User `json:"users"` // Using user.User directly from the service layer
	TotalCount int         `json:"total_count"`
}

// makeListUsersEndpoint creates an endpoint for listing users.
func makeListUsersEndpoint(s user.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Parse query parameters for pagination
		offsetStr := c.DefaultQuery("offset", "0")
		limitStr := c.DefaultQuery("limit", "10")

		offset, _ := strconv.Atoi(offsetStr)
		limit, _ := strconv.Atoi(limitStr)

		// Create a query object, assuming pkg/utl/query has a constructor or can be initialized like this
		q := &query.ListQuery{
			Offset: offset,
			Limit:  limit,
			Sort:   c.DefaultQuery("sort", ""), // Assuming sort parameter as well
		}

		userList, err := s.ListUsers(c.Request.Context(), q)
		if err != nil {
			// A proper error handling middleware or helper would be used here.
			// For simplicity, directly returning error.
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, ListResponse{Users: userList.Users, TotalCount: userList.TotalCount})
	}
}

// Assuming other endpoints like makeCreateUserEndpoint, makeViewUserEndpoint etc. exist
// ...

// RegisterHandlers registers the HTTP handlers for the user service.
func RegisterHandlers(r *gin.Engine, s user.Service) {
	users := r.Group("/users")
	// Add any necessary middleware for the group, e.g., authentication, RBAC
	// users.Use(middleware.AuthMiddleware(), middleware.RBACMiddleware())

	users.GET("", makeListUsersEndpoint(s))
	// users.POST("", makeCreateUserEndpoint(s))
	// users.GET("/:id", makeViewUserEndpoint(s))
	// users.DELETE("/:id", makeDeleteUserEndpoint(s))
}
