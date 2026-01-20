package handlers

import (
	"github.com/google/uuid"

	"github.com/yourusername/go-starter/internal/repository"
)

// UserResponse represents the user data returned in API responses
// Contains only minimal fields: ID, Name, Email
type UserResponse struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// NewUserResponse creates a UserResponse from a repository User model
func NewUserResponse(user *repository.User) *UserResponse {
	return &UserResponse{
		Name:  user.Name,
		Email: user.Email,
	}
}

// ToJSONAPIData converts a user to JSON:API data format
func ToJSONAPIData(user *repository.User) JSONAPIData {
	return JSONAPIData{
		Type:       "users",
		ID:         uuid.UUID(user.ID).String(),
		Attributes: NewUserResponse(user),
	}
}
