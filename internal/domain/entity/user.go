// Package entity contains pure domain entities — no ORM or transport tags.
package entity

// User is the immutable domain entity. Persistence-layer ORM model
// (infrastructure/persistence/postgres/model.User) is kept separate
// to preserve dependency direction: Domain never knows about storage.
type User struct {
	LastName  string
	FirstName string
	Email     string
	// Password is plain-text at the domain boundary; hashing is delegated
	// to the domain service so the entity itself never carries credentials
	// past UserCreator.
	Password string
}
