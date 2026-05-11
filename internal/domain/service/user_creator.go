// Package service hosts domain services — pure orchestration of entities,
// VOs and repository contracts. No HTTP, no SQL, no Kafka here.
package service

import (
	"context"
	"errors"
	"fmt"

	"api/internal/domain/entity"
	"api/internal/domain/event"
	"api/internal/domain/repository"
)

// ErrUserNotPersisted is returned when the repository returns a nil id —
// meaning the row was rejected for a reason caller couldn't predict.
var ErrUserNotPersisted = errors.New("user was not persisted")

// UserCreator orchestrates registration: hash credentials, persist, then
// publish a UserRegistered event so async consumers can react.
type UserCreator struct {
	users    repository.UserRepository
	hasher   PasswordHasher
	eventBus event.Bus
}

// NewUserCreator wires the collaborators. All three are required.
func NewUserCreator(
	users repository.UserRepository,
	hasher PasswordHasher,
	eventBus event.Bus,
) *UserCreator {
	return &UserCreator{users: users, hasher: hasher, eventBus: eventBus}
}

// Create hashes the user's password, stores the entity, and publishes a
// UserRegistered event. Event-publish failures are returned to the caller
// rather than swallowed: the original PHP project relies on
// transactional outbox / retry at a higher level. If you later add an
// outbox, replace the direct Publish here.
func (s *UserCreator) Create(ctx context.Context, user entity.User) (int, error) {
	hash, err := s.hasher.Hash(user.Password)
	if err != nil {
		return 0, fmt.Errorf("hash password: %w", err)
	}

	idPtr, err := s.users.Store(ctx, user, hash)
	if err != nil {
		return 0, fmt.Errorf("store user: %w", err)
	}
	if idPtr == nil {
		return 0, ErrUserNotPersisted
	}

	if err := s.eventBus.Publish(event.UserRegistered{ID: *idPtr}); err != nil {
		return 0, fmt.Errorf("publish user_registered: %w", err)
	}
	return *idPtr, nil
}
