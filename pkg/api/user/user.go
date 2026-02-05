package user

import (
	"context"

	"github.com/user/repo/pkg/utl/query" // Assuming this path for query package
)

// User represents a user. (Assuming this struct exists as part of the domain model)
type User struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	// ... other fields relevant to a user ...
}

// UserList represents a list of users with total count for pagination.
type UserList struct {
	Users      []User `json:"users"`
	TotalCount int    `json:"total_count"`
}

// Service specifies the user service interface.
type Service interface {
	Create(ctx context.Context, user User) (User, error)
	View(ctx context.Context, id int) (User, error)
	Delete(ctx context.Context, id int) error
	ListUsers(ctx context.Context, q *query.ListQuery) (UserList, error) // Modified return type
	// ... other methods ...
}

// Repository specifies the user repository interface.
type Repository interface {
	Create(ctx context.Context, user User) (User, error)
	View(ctx context.Context, id int) (User, error)
	Delete(ctx context.Context, id int) error
	List(ctx context.Context, q *query.ListQuery) ([]User, int, error) // Modified return type to include count
	// ... other methods ...
}

// service implements Service interface
type service struct {
	repo Repository
	// ... other dependencies like Logger, etc.
}

// NewService creates a new user service.
func NewService(repo Repository) Service {
	return &service{repo: repo}
}

// ListUsers returns a list of users.
func (s *service) ListUsers(ctx context.Context, q *query.ListQuery) (UserList, error) {
	users, totalCount, err := s.repo.List(ctx, q) // Call updated repo method
	if err != nil {
		return UserList{}, err
	}
	return UserList{Users: users, TotalCount: totalCount}, nil
}

// ... other service method implementations would follow ...
// Example:
// func (s *service) Create(ctx context.Context, user User) (User, error) { /* ... */ }
// func (s *service) View(ctx context.Context, id int) (User, error) { /* ... */ }
// func (s *service) Delete(ctx context.Context, id int) error { /* ... */ }
