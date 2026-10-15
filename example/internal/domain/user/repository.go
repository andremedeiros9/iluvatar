package user

import (
	"context"

	"github.com/google/uuid"
)

// Repository persists User records.
type Repository interface {
	// Create inserts m, setting its ID, CreatedAt, and UpdatedAt.
	Create(ctx context.Context, m *User) error
	// Get returns the user with the given id.
	Get(ctx context.Context, id uuid.UUID) (*User, error)
	// List returns every user.
	List(ctx context.Context) ([]User, error)
	// Update overwrites the stored fields of m, refreshing UpdatedAt.
	Update(ctx context.Context, m *User) error
	// Delete removes the user with the given id.
	Delete(ctx context.Context, id uuid.UUID) error
}
