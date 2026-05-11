// Package repository declares the persistence contracts owned by the
// domain. Concrete implementations live in infrastructure.
package repository

import (
	"context"

	"api/internal/domain/entity"
)

// UserRepository persists User aggregates. Returning *int matches the
// original PHP signature: nil means "store did not produce an id" which
// the caller must treat as an error.
type UserRepository interface {
	Store(ctx context.Context, user entity.User, hashPassword string) (*int, error)
}
