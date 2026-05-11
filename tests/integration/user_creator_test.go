package integration_test

import (
	"context"
	"errors"
	"testing"

	"api/internal/domain/entity"
	"api/internal/domain/event"
	"api/internal/domain/service"

	"github.com/stretchr/testify/assert"
)

// stubHasher returns the input as-is, prefixed so we can assert it
// arrived at the repository.
type stubHasher struct{}

func (stubHasher) Hash(p string) (string, error)   { return "hashed:" + p, nil }
func (stubHasher) Verify(_, _ string) error        { return nil }

// nilStoringRepo simulates a repository that "stores" without producing
// an id — same path as the PHP integration test, where Store returns
// null and the service must surface a hard error.
type nilStoringRepo struct {
	called      bool
	savedHash   string
	savedEntity entity.User
}

func (r *nilStoringRepo) Store(_ context.Context, u entity.User, hash string) (*int, error) {
	r.called = true
	r.savedEntity = u
	r.savedHash = hash
	return nil, nil
}

// inMemoryBus captures publishes for assertion.
type inMemoryBus struct{ events []event.DomainEvent }

func (b *inMemoryBus) Publish(e event.DomainEvent) error {
	b.events = append(b.events, e)
	return nil
}

func TestUserCreator_StoreReturnsNilID_ProducesError(t *testing.T) {
	repo := &nilStoringRepo{}
	bus := &inMemoryBus{}
	svc := service.NewUserCreator(repo, stubHasher{}, bus)

	_, err := svc.Create(context.Background(), entity.User{
		FirstName: "Ada", LastName: "Lovelace",
		Email: "ada@example.com", Password: "secret",
	})

	assert.True(t, repo.called, "repository should be called")
	assert.Equal(t, "hashed:secret", repo.savedHash, "hashed password should be forwarded")
	assert.ErrorIs(t, err, service.ErrUserNotPersisted)
	assert.Empty(t, bus.events, "no event should be published when persist fails")
}

// repoOK returns a fixed id.
type repoOK struct{}

func (repoOK) Store(_ context.Context, _ entity.User, _ string) (*int, error) {
	id := 42
	return &id, nil
}

// publishErrBus simulates kafka being unreachable.
type publishErrBus struct{}

func (publishErrBus) Publish(_ event.DomainEvent) error { return errors.New("broker down") }

func TestUserCreator_PublishFailure_PropagatesError(t *testing.T) {
	svc := service.NewUserCreator(repoOK{}, stubHasher{}, publishErrBus{})
	_, err := svc.Create(context.Background(), entity.User{
		Email: "a@b.c", Password: "p",
	})
	assert.Error(t, err)
}
