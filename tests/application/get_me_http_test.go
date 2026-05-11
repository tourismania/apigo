package application_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"api/internal/application/bus"
	"api/internal/infrastructure/auth"
	getmehttp "api/internal/presentation/http/api/v1/get_me"
	custommw "api/internal/presentation/http/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

// Mirrors the PHP application test: GET /api/v1/me without a token must
// return 401. We mount the same JWT middleware the production router uses
// to keep this an end-to-end-ish assertion.
func TestGetMe_RequiresAuth(t *testing.T) {
	// jwt.Service is intentionally constructed with a nil key — the middleware
	// rejects the request before any key material is touched because the
	// Authorization header is absent.
	var jwtSvc *auth.Service // nil is fine; never reached.

	queryBus := bus.NewInMemoryQueryBus()
	handler := getmehttp.NewHandler(queryBus, getmehttp.NewResolver())

	r := chi.NewRouter()
	r.Route("/api/v1", func(api chi.Router) {
		api.Use(custommw.JWT(jwtSvc))
		api.Get("/me", handler.Handle)
	})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/me", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}
