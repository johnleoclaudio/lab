package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/yourusername/go-starter/internal/db"
	"github.com/yourusername/go-starter/internal/models"
)

// User represents the domain model for a user
type User struct {
	ID        uuid.UUID
	Email     string
	Name      string
	CreatedAt pgtype.Timestamptz
	UpdatedAt pgtype.Timestamptz
}

// UserRepository defines the interface for user data access
type UserRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*User, error)
}

// userRepository implements UserRepository
type userRepository struct {
	queries *db.Queries
}

// NewUserRepository creates a new UserRepository
func NewUserRepository(queries *db.Queries) UserRepository {
	return &userRepository{
		queries: queries,
	}
}

// GetByID retrieves a user by their ID
func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*User, error) {
	// Convert uuid.UUID to pgtype.UUID
	pgID := pgtype.UUID{
		Bytes: id,
		Valid: true,
	}

	// Query the database using sqlc-generated code
	dbUser, err := r.queries.GetUserByID(ctx, pgID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNotFound
		}
		return nil, fmt.Errorf("get user by id: %w", err)
	}

	// Convert database model to domain model
	user := &User{
		ID:        uuid.UUID(dbUser.ID.Bytes),
		Email:     dbUser.Email,
		Name:      dbUser.Name,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
	}

	return user, nil
}
