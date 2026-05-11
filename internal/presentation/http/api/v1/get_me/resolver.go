package getmehttp

import (
	"context"
	"errors"

	"api/internal/presentation/http/middleware"
)

// ErrNoAuthClaims means the request didn't traverse JWT middleware. This
// is a 500-tier programmer error, not a 401 — the route is misconfigured.
var ErrNoAuthClaims = errors.New("no auth claims on context")

// ErrUserMissingID covers the analogous Symfony path: a token without an
// identifier reaches us, which we report as a server-side problem
// because the issuer is supposed to enforce that.
var ErrUserMissingID = errors.New("token has no user id")

// Resolver reads the authenticated principal off the request context.
// Kept as a separate type (rather than inline in the handler) for parity
// with the PHP GetMeResolver and so it can be unit-tested in isolation.
type Resolver struct{}

// NewResolver constructs the resolver.
func NewResolver() *Resolver { return &Resolver{} }

// Resolve returns the DTO for the authenticated user, or an error
// describing why it can't.
func (r *Resolver) Resolve(ctx context.Context) (*GetMeDto, error) {
	claims, ok := middleware.ClaimsFromContext(ctx)
	if !ok || claims == nil {
		return nil, ErrNoAuthClaims
	}
	if claims.UID == 0 {
		return nil, ErrUserMissingID
	}
	return &GetMeDto{
		ID:        claims.UID,
		Email:     emailOrUsername(claims.Email, claims.Username),
		FirstName: "",
		LastName:  "",
		Phone:     "",
		Roles:     append([]string(nil), claims.Roles...),
	}, nil
}

func emailOrUsername(email, username string) string {
	if email != "" {
		return email
	}
	return username
}
