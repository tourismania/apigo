package createuser

import (
	"context"

	"api/internal/domain/entity"
	"api/internal/domain/service"
)

// Handler executes the CreateUser command by delegating to the domain
// UserCreator service. Keeping the handler thin preserves DDD: business
// invariants stay in the domain layer.
type Handler struct {
	userCreator *service.UserCreator
}

// NewHandler constructs the handler.
func NewHandler(userCreator *service.UserCreator) *Handler {
	return &Handler{userCreator: userCreator}
}

// Handle implements the typed handler contract registered on the
// CommandBus.
func (h *Handler) Handle(ctx context.Context, cmd Command) (Result, error) {
	id, err := h.userCreator.Create(ctx, entity.User{
		FirstName: cmd.FirstName,
		LastName:  cmd.LastName,
		Email:     cmd.Email,
		Password:  cmd.Password,
	})
	if err != nil {
		return Result{}, err
	}
	return Result{ID: id}, nil
}
