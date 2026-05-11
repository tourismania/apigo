package getmehttp

import (
	"errors"
	"net/http"

	"api/internal/application/bus"
	getmeq "api/internal/application/query/get_me"
	"api/internal/presentation/http/httpx"
)

// Handler renders the authenticated user as a JSON document. It is
// thin — the heavy lifting (rights derivation) lives in the query
// handler. The handler just maps DTO ↔ query ↔ response.
type Handler struct {
	bus      bus.QueryBus
	resolver *Resolver
}

// NewHandler constructs the handler.
func NewHandler(b bus.QueryBus, resolver *Resolver) *Handler {
	return &Handler{bus: b, resolver: resolver}
}

// Handle is the http.HandlerFunc.
//
//	@Summary      Current user profile
//	@Description  Returns the authenticated user's profile and computed rights.
//	@Tags         Profile
//	@Produce      json
//	@Success      200  {object}  GetMeResponse
//	@Failure      401  {object}  httpx.ErrorBody
//	@Failure      404  {object}  httpx.ErrorBody
//	@Failure      500  {object}  httpx.ErrorBody
//	@Security     BearerAuth
//	@Router       /api/v1/me [get]
func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	dto, err := h.resolver.Resolve(r.Context())
	if err != nil {
		switch {
		case errors.Is(err, ErrNoAuthClaims):
			httpx.WriteError(w, http.StatusUnauthorized, "unauthenticated")
		case errors.Is(err, ErrUserMissingID):
			httpx.WriteError(w, http.StatusInternalServerError, "token missing user id")
		default:
			httpx.WriteError(w, http.StatusNotFound, "user not found")
		}
		return
	}

	raw, err := h.bus.Dispatch(r.Context(), getmeq.Query{
		ID:        dto.ID,
		Email:     dto.Email,
		Phone:     dto.Phone,
		FirstName: dto.FirstName,
		LastName:  dto.LastName,
		Roles:     dto.Roles,
	})
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	res, ok := raw.(getmeq.Result)
	if !ok {
		httpx.WriteError(w, http.StatusInternalServerError, "unexpected handler result type")
		return
	}

	httpx.WriteJSON(w, http.StatusOK, GetMeResponse{
		ID:        res.ID,
		Email:     res.Email,
		Phone:     res.Phone,
		FirstName: res.FirstName,
		LastName:  res.LastName,
		Rights:    Rights{IsSuperAdmin: res.Rights.IsSuperAdmin},
	})
}
