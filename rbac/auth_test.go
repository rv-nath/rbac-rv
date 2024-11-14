package rbac

import (
	"testing"
)

func TestCheckPermission(t *testing.T) {
	// Mock functions for user roles and role permissions
	fetchUserRoles := func(userID string) (UserRoles, error) {
		return UserRoles{Roles: []string{"admin"}}, nil
	}
	fetchRolePermissions := func(roleID string) (RolePermissions, error) {
		return RolePermissions{
			Perms: []Permission{
				{Resource: "user", Permission: "view", Allowed: true, Scope: "self"},
				{Resource: "user", Permission: "edit", Allowed: false, Scope: "self"},
				{Resource: "order", Permission: "create", Allowed: true, Scope: "global"},
			},
		}, nil
	}

	rbac := NewRBAC(fetchUserRoles, fetchRolePermissions)

	tests := []struct {
		name           string
		userID         string
		resourceID     string
		resource       string
		action         string
		scope          string
		expectedResult bool
	}{
		{"Self Scope View Allowed", "user1", "123", "user", "view", "self", true},
		{"Self Scope Edit Disallowed", "user1", "123", "user", "edit", "self", false},
		{"Global Scope Create Allowed", "user1", "", "order", "create", "global", true},
		{"Unauthorized Action", "user1", "123", "order", "edit", "self", false},
	}

	for _, tt := range tests {
		t.Logf("Running test %s...", tt.name)
		t.Run(tt.name, func(t *testing.T) {
			result := rbac.CheckPermission(tt.userID, tt.resourceID, tt.resource, tt.action, tt.scope)
			if result != tt.expectedResult {
				t.Errorf("Test %s failed: got %v, want %v", tt.name, result, tt.expectedResult)
			}
		})
	}
}
