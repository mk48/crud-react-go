package model

import (
	"time"

	"github.com/google/uuid"
)

// User mirrors the "user" table. It is shared across services (rather than
// living in a "user" service package) because the auth middleware needs to
// create/look it up before any service package exists to own it.
type User struct {
	ID        uuid.UUID  `db:"id"`
	Sub       string     `db:"sub"`
	Email     string     `db:"email"`
	Name      *string    `db:"name"`
	IsAdmin   bool       `db:"is_admin"`
	CreatedAt time.Time  `db:"created_at"`
	CreatedBy uuid.UUID  `db:"created_by"`
	UpdatedAt *time.Time `db:"updated_at"`
	UpdatedBy *uuid.UUID `db:"updated_by"`
	DeletedAt *time.Time `db:"deleted_at"`
	DeletedBy *uuid.UUID `db:"deleted_by"`
}
