// Package bus implements thin in-memory command and query buses, mirroring
// the Symfony Messenger pattern: each message type is associated with a
// single handler resolved by reflect type at dispatch time.
package bus

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"sync"
)

// CommandHandlerFunc is the canonical handler signature stored in the
// registry. We hide the generic CommandHandler[C,R] behind it so the bus
// can carry handlers for heterogeneous command types in one map.
type CommandHandlerFunc func(ctx context.Context, cmd any) (any, error)

// CommandBus dispatches a command to its registered handler.
type CommandBus interface {
	Dispatch(ctx context.Context, cmd any) (any, error)
}

// CommandRegistry exposes registration; kept separate so application code
// can require only the read-side of the bus.
type CommandRegistry interface {
	Register(commandType reflect.Type, handler CommandHandlerFunc)
}

// InMemoryCommandBus is the default routing implementation.
type InMemoryCommandBus struct {
	mu       sync.RWMutex
	handlers map[reflect.Type]CommandHandlerFunc
}

// NewInMemoryCommandBus constructs an empty bus.
func NewInMemoryCommandBus() *InMemoryCommandBus {
	return &InMemoryCommandBus{handlers: make(map[reflect.Type]CommandHandlerFunc)}
}

// Register associates a handler with a command type. Re-registration for
// the same type panics: dispatching to an ambiguous handler is a
// programmer error, not a runtime condition.
func (b *InMemoryCommandBus) Register(commandType reflect.Type, handler CommandHandlerFunc) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if _, exists := b.handlers[commandType]; exists {
		panic(fmt.Sprintf("command handler already registered for %s", commandType))
	}
	b.handlers[commandType] = handler
}

// ErrNoHandler is returned when Dispatch sees an unknown command.
var ErrNoHandler = errors.New("no handler registered for command")

// Dispatch routes the command to its registered handler.
func (b *InMemoryCommandBus) Dispatch(ctx context.Context, cmd any) (any, error) {
	if cmd == nil {
		return nil, errors.New("nil command")
	}
	t := reflect.TypeOf(cmd)
	b.mu.RLock()
	handler, ok := b.handlers[t]
	b.mu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrNoHandler, t)
	}
	return handler(ctx, cmd)
}

// RegisterCommand is a typed convenience wrapper around Register: it
// adapts a strongly typed handler into the untyped registry entry.
func RegisterCommand[C any, R any](
	registry CommandRegistry,
	handler func(ctx context.Context, cmd C) (R, error),
) {
	var zero C
	registry.Register(reflect.TypeOf(zero), func(ctx context.Context, raw any) (any, error) {
		typed, ok := raw.(C)
		if !ok {
			return nil, fmt.Errorf("command bus: expected %T, got %T", zero, raw)
		}
		return handler(ctx, typed)
	})
}
