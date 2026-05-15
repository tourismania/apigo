// Package getmehttp is the HTTP boundary for the GetMe query.
package getmehttp

// GetMeDto is the transport view of the authenticated user, populated by
// the resolver from JWT claims.
type GetMeDto struct {
	ID        int
	Email     string
	Phone     string
	FirstName string
	LastName  string
	Roles     []string
}

// GetMeResponse is what we serialise back to the client. Rights is
// flattened explicitly to keep the wire payload stable even if the
// underlying Value Object grows new fields.
type GetMeResponse struct {
	ID        int    `json:"id"`
	Email     string `json:"email"`
	Phone     string `json:"phone,omitempty"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Rights    Rights `json:"rights"`
}

// Rights is the public projection of valueobject.RightsDescribe.
type Rights struct {
	IsSuperAdmin bool `json:"is_super_admin"`
}
