// Package loginhttp is the HTTP boundary for the password-grant login.
package loginhttp

// LoginRequest mirrors lexik_jwt's check_path body: {"username","password"}.
type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// LoginResponse is the issued token envelope.
type LoginResponse struct {
	Token string `json:"token"`
}
