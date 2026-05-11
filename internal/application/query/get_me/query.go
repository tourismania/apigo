// Package getme contains the "current user profile" read-side use case.
package getme

// Query carries everything we already know about the authenticated user
// (resolved from JWT claims at the presentation boundary). The handler
// then enriches it with derived data — e.g. RightsDescribe.
type Query struct {
	ID        int
	Email     string
	Phone     string
	FirstName string
	LastName  string
	Roles     []string
}
