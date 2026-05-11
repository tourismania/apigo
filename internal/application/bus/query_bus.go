package bus

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"sync"
)

// QueryHandlerFunc has the same shape as CommandHandlerFunc but is kept
// distinct so CQRS responsibilities don't get crossed at dispatch time.
type QueryHandlerFunc func(ctx context.Context, query any) (any, error)

// QueryBus dispatches a query to its registered handler.
type QueryBus interface {
	Dispatch(ctx context.Context, query any) (any, error)
}

// QueryRegistry exposes registration.
type QueryRegistry interface {
	Register(queryType reflect.Type, handler QueryHandlerFunc)
}

// InMemoryQueryBus is the default routing implementation.
type InMemoryQueryBus struct {
	mu       sync.RWMutex
	handlers map[reflect.Type]QueryHandlerFunc
}

// NewInMemoryQueryBus constructs an empty bus.
func NewInMemoryQueryBus() *InMemoryQueryBus {
	return &InMemoryQueryBus{handlers: make(map[reflect.Type]QueryHandlerFunc)}
}

// Register associates a handler with a query type. Duplicate registration
// for the same type is treated as a programmer error.
func (b *InMemoryQueryBus) Register(queryType reflect.Type, handler QueryHandlerFunc) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if _, exists := b.handlers[queryType]; exists {
		panic(fmt.Sprintf("query handler already registered for %s", queryType))
	}
	b.handlers[queryType] = handler
}

// Dispatch routes the query to its registered handler.
func (b *InMemoryQueryBus) Dispatch(ctx context.Context, query any) (any, error) {
	if query == nil {
		return nil, errors.New("nil query")
	}
	t := reflect.TypeOf(query)
	b.mu.RLock()
	handler, ok := b.handlers[t]
	b.mu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrNoHandler, t)
	}
	return handler(ctx, query)
}

// RegisterQuery is the typed convenience wrapper, analogous to
// RegisterCommand.
func RegisterQuery[Q any, R any](
	registry QueryRegistry,
	handler func(ctx context.Context, q Q) (R, error),
) {
	var zero Q
	registry.Register(reflect.TypeOf(zero), func(ctx context.Context, raw any) (any, error) {
		typed, ok := raw.(Q)
		if !ok {
			return nil, fmt.Errorf("query bus: expected %T, got %T", zero, raw)
		}
		return handler(ctx, typed)
	})
}
