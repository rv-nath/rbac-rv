package rbac

import (
	"testing"
)

func TestDetermineIntent(t *testing.T) {
	fetchResources := func() ([]string, error) {
		return []string{"user", "order", "product"}, nil
	}

	tests := []struct {
		name               string
		method             string
		path               string
		expectedAction     string
		expectedResource   string
		expectedResourceID string
		shouldError        bool
	}{
		{"View Users List", "GET", "/api/v1/mod1/users", "view", "user", "", false},
		{"Create User", "POST", "/api/v1/mod1/users", "create", "user", "", false},
		{"View Specific User", "GET", "/api/v1/mod1/users/123", "view", "user", "123", false},
		{"View Specific Order", "GET", "/api/v1/mod1/orders/757f45f8-91f5-4c93-bae8-04b0d4b2a75b", "view", "order", "757f45f8-91f5-4c93-bae8-04b0d4b2a75b", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("Running test: %s", tt.name)
			action, resource, resourceID, err := DetermineIntent(tt.method, tt.path, fetchResources)
			if tt.shouldError && err == nil {
				t.Errorf("Expected error but got none")
			} else if !tt.shouldError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			} else if action != tt.expectedAction || resource != tt.expectedResource || resourceID != tt.expectedResourceID {
				t.Errorf("Expected action=%s, resource=%s, resourceID=%s; got action=%s, resource=%s, resourceID=%s",
					tt.expectedAction, tt.expectedResource, tt.expectedResourceID, action, resource, resourceID)
			}
			t.Logf("Finished test: %s", tt.name)
		})
	}
}
