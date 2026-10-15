package user

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// Repository stores User records in user.
type Repository struct {
	db *sqlx.DB
}

// NewRepository returns a Repository backed by db.
func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

// Create inserts m, setting its ID, CreatedAt, and UpdatedAt.
func (r *Repository) Create(ctx context.Context, m *User) error {
	m.ID = uuid.New()
	now := time.Now().UTC()
	m.CreatedAt = now
	m.UpdatedAt = now

	query := r.db.Rebind(`INSERT INTO user (id, email, age, created_at, updated_at) VALUES (?, ?, ?, ?, ?)`)
	_, err := r.db.ExecContext(ctx, query,
		m.ID,
		m.Email,
		m.Age,
		m.CreatedAt,
		m.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("creating user: %w", err)
	}
	return nil
}

// Get returns the user with the given id.
func (r *Repository) Get(ctx context.Context, id uuid.UUID) (*User, error) {
	var m User
	query := r.db.Rebind(`SELECT id, email, age, created_at, updated_at FROM user WHERE id = ?`)
	if err := r.db.GetContext(ctx, &m, query, id); err != nil {
		return nil, fmt.Errorf("getting user %s: %w", id, err)
	}
	return &m, nil
}

// List returns every user.
func (r *Repository) List(ctx context.Context) ([]User, error) {
	var items []User
	query := r.db.Rebind(`SELECT id, email, age, created_at, updated_at FROM user ORDER BY created_at`)
	if err := r.db.SelectContext(ctx, &items, query); err != nil {
		return nil, fmt.Errorf("listing users: %w", err)
	}
	return items, nil
}

// Update overwrites the stored fields of m, refreshing UpdatedAt.
func (r *Repository) Update(ctx context.Context, m *User) error {
	m.UpdatedAt = time.Now().UTC()

	query := r.db.Rebind(`UPDATE user SET email = ?, age = ?, updated_at = ? WHERE id = ?`)
	_, err := r.db.ExecContext(ctx, query,
		m.Email,
		m.Age,
		m.UpdatedAt,
		m.ID,
	)
	if err != nil {
		return fmt.Errorf("updating user %s: %w", m.ID, err)
	}
	return nil
}

// Delete removes the user with the given id.
func (r *Repository) Delete(ctx context.Context, id uuid.UUID) error {
	query := r.db.Rebind(`DELETE FROM user WHERE id = ?`)
	if _, err := r.db.ExecContext(ctx, query, id); err != nil {
		return fmt.Errorf("deleting user %s: %w", id, err)
	}
	return nil
}
