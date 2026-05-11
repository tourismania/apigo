// Package kafka contains the segmentio-kafka-go adapter that publishes
// domain events.
package kafka

import (
	"encoding/json"
	"errors"
)

// Encoder is the local view of a domain event the serializer cares about.
// We re-declare it (instead of importing event.DomainEvent) so the kafka
// package depends only on this minimal contract.
type Encoder interface {
	GetKey() string
	GetEventCode() string
}

// Encode serialises a domain event into the canonical kafka payload:
//
//	key  = event.GetKey()
//	body = JSON({...event fields, "code": event.GetEventCode()})
//
// We marshal the event into a generic map first so the "code"
// discriminator is appended without requiring every event to embed it
// explicitly.
func Encode(event Encoder) (string, []byte, error) {
	if event == nil {
		return "", nil, errors.New("nil event")
	}

	raw, err := json.Marshal(event)
	if err != nil {
		return "", nil, err
	}
	var bag map[string]any
	if err := json.Unmarshal(raw, &bag); err != nil {
		return "", nil, err
	}
	if bag == nil {
		bag = map[string]any{}
	}
	bag["code"] = event.GetEventCode()

	body, err := json.Marshal(bag)
	if err != nil {
		return "", nil, err
	}
	return event.GetKey(), body, nil
}
