package unit_test

import (
	"testing"

	"api/internal/domain/enum"
	"api/internal/domain/factory"

	"github.com/stretchr/testify/assert"
)

// Mirrors the original PHP unit test: roles → IsSuperAdmin truth table.
func TestRightsDescribeFactory_ByRoles(t *testing.T) {
	f := factory.NewRightsDescribeFactory()

	cases := []struct {
		name  string
		roles []string
		want  bool
	}{
		{name: "empty", roles: []string{}, want: false},
		{name: "only super admin", roles: []string{string(enum.RoleSuperAdmin)}, want: true},
		{name: "only user", roles: []string{string(enum.RoleUser)}, want: false},
		{name: "mixed includes super", roles: []string{string(enum.RoleUser), string(enum.RoleSuperAdmin)}, want: true},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got := f.ByRoles(tc.roles)
			assert.Equal(t, tc.want, got.IsSuperAdmin)
		})
	}
}
