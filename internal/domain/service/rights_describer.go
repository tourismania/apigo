package service

import (
	"api/internal/domain/factory"
	"api/internal/domain/valueobject"
)

// RightsDescriber is a thin domain service that wraps the factory so
// consumers (e.g. GetMe handler) don't need to instantiate the factory
// themselves and the dependency stays explicit in DI.
type RightsDescriber struct {
	factory *factory.RightsDescribeFactory
}

// NewRightsDescriber constructs the service.
func NewRightsDescriber(f *factory.RightsDescribeFactory) *RightsDescriber {
	return &RightsDescriber{factory: f}
}

// ByRoles delegates to the underlying factory; kept on a separate type so
// future audit/logging hooks can attach without changing the factory.
func (s *RightsDescriber) ByRoles(roles []string) valueobject.RightsDescribe {
	return s.factory.ByRoles(roles)
}
