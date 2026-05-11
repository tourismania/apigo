package db

import (
	"time"

	"github.com/google/uuid"
)

// User is the sqlc-style row representation. Nullable columns use the
// pgx-friendly pointer form; JSON column is decoded into a free-form
// map. Mapping to the domain entity happens in the repository adapter.
type User struct {
	ID               int32
	Uuid             uuid.UUID
	FirstName        *string
	LastName         *string
	Email            string
	Login            string
	CreatedAt        time.Time
	UpdatedAt        time.Time
	Phone            *string
	Password         string
	IsActive         bool
	Birthday         *time.Time
	ExtraInformation []byte
	Roles            []string
}
