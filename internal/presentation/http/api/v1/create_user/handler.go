package createuserhttp

import (
	"errors"
	"net/http"

	"api/internal/application/bus"
	createusercmd "api/internal/application/command/create_user"
	"api/internal/presentation/http/httpx"

	"github.com/go-playground/validator/v10"
)

// Handler turns HTTP requests into CreateUser commands dispatched to the
// command bus. Validation, decoding and response shaping live here so
// the bus / handler stay transport-agnostic.
type Handler struct {
	bus      bus.CommandBus
	validate *validator.Validate
}

// NewHandler constructs the handler.
func NewHandler(b bus.CommandBus, v *validator.Validate) *Handler {
	return &Handler{bus: b, validate: v}
}

// Handle is the http.HandlerFunc.
//
//	@Summary      Create a user
//	@Description  Registers a new user account.
//	@Tags         Users
//	@Accept       json
//	@Produce      json
//	@Param        body  body      CreateUserRequest  true  "User payload"
//	@Success      201   {object}  CreateUserResponse
//	@Failure      400   {object}  httpx.ErrorBody
//	@Failure      401   {object}  httpx.ErrorBody
//	@Failure      500   {object}  httpx.ErrorBody
//	@Security     BearerAuth
//	@Router       /api/v1/users [post]
func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest
	if err := httpx.DecodeJSON(r, &req, h.validate); err != nil {
		if errors.Is(err, httpx.ErrBadJSON) {
			httpx.WriteError(w, http.StatusBadRequest, "invalid JSON body")
			return
		}
		httpx.WriteValidationError(w, err)
		return
	}

	raw, err := h.bus.Dispatch(r.Context(), createusercmd.Command{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Password:  req.Password,
	})
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	res, ok := raw.(createusercmd.Result)
	if !ok {
		httpx.WriteError(w, http.StatusInternalServerError, "unexpected handler result type")
		return
	}
	httpx.WriteJSON(w, http.StatusCreated, CreateUserResponse{ID: res.ID})
}
