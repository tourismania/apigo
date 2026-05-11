package event

import "strconv"

// UserRegistered fires after a user is persisted. The body is intentionally
// minimal — downstream services should rehydrate from the DB by ID.
type UserRegistered struct {
	ID int `json:"id"`
}

const userRegisteredCode = "user_registered"

// GetKey is used as the Kafka partition key.
func (e UserRegistered) GetKey() string { return strconv.Itoa(e.ID) }

// GetEventCode is the discriminator written into the message body.
func (e UserRegistered) GetEventCode() string { return userRegisteredCode }
