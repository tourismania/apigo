// Package model holds persistence-only structs. They mirror the database
// schema and are intentionally separate from domain entities.
package model

import (
	"time"

	"github.com/google/uuid"
)

// User is the ORM/row representation. Nullable columns use pointers so
// scanning won't fail on NULLs. ExtraInformation is decoded from the
// JSON column into a free-form map; if you later add a strict schema,
// replace with a typed struct.
type User struct {
	ID               int
	UUID             uuid.UUID
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
	ExtraInformation map[string]any
	Roles            []string
}
