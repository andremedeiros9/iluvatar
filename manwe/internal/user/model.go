package user

import (
	"time"

	"github.com/google/uuid"
)

// User is a user record.
type User struct {
	ID        uuid.UUID `db:"id"`
	Email     string    `db:"email"`
	Age       int64     `db:"age"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
