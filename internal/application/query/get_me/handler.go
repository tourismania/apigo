package getme

import (
	"context"

	"api/internal/domain/service"
)

// UseCase is the port the presentation layer depends on.
type UseCase interface {
	Handle(ctx context.Context, q Query) (Result, error)
}

// Handler derives the rights bag from the authenticated user's roles and
// echoes the rest of the profile. It performs no I/O — everything it
// needs is already in the Query carrier supplied by the resolver.
type Handler struct {
	rightsDescriber *service.RightsDescriber
}

// NewHandler constructs the handler.
func NewHandler(rightsDescriber *service.RightsDescriber) *Handler {
	return &Handler{rightsDescriber: rightsDescriber}
}

// Handle satisfies UseCase.
func (h *Handler) Handle(_ context.Context, q Query) (Result, error) {
	rights := h.rightsDescriber.ByRoles(q.Roles)
	return Result{
		ID:        q.ID,
		Email:     q.Email,
		Phone:     q.Phone,
		FirstName: q.FirstName,
		LastName:  q.LastName,
		Rights:    rights,
	}, nil
}
