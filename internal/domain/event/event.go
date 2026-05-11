// Package event defines domain events and the publishing contract.
package event

// DomainEvent is the marker interface for any business event emitted by
// the domain. Serializers in infrastructure consume it via the Encoder
// view (GetKey/GetEventCode) — domain itself stays unaware of the wire
// format.
type DomainEvent interface {
	GetKey() string
	GetEventCode() string
}

// Bus publishes domain events to subscribers (kafka in infrastructure).
// Kept here so domain services depend only on this contract.
type Bus interface {
	Publish(event DomainEvent) error
}
