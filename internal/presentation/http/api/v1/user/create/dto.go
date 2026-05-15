// Package createuserhttp is the HTTP boundary for the CreateUser
// command — request decoding, validation and response shaping.
package createuserhttp

// CreateUserRequest is the public request body schema.
type CreateUserRequest struct {
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name"  validate:"required"`
	Email     string `json:"email"      validate:"required,email"`
	Password  string `json:"password"   validate:"required,min=6"`
}

// CreateUserResponse is the response envelope.
type CreateUserResponse struct {
	ID int `json:"id"`
}
