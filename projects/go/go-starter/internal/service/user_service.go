package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/yourusername/go-starter/internal/repository"
)

// UserService defines the interface for user business logic
type UserService interface {
	GetUser(ctx context.Context, id uuid.UUID) (*repository.User, error)
}

// userService implements UserService
type userService struct {
	userRepo repository.UserRepository
}

// NewUserService creates a new UserService
func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

// GetUser retrieves a user by their ID
func (s *userService) GetUser(ctx context.Context, id uuid.UUID) (*repository.User, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}

	return user, nil
}
