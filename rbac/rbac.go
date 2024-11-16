package rbac

import (
	"sync"
	// log "github.com/sirupsen/logrus"
)

// Userroles stores the list of roles associated with a user.
type UserRoles struct {
	Roles []string // List of roles for mapped to a user.
}

// Permission defines a permission for a resource in a given scope
type Permission struct {
	Resource   string
	Permission string
	Allowed    bool
	Scope      string
}

// RolePermissions stores permissions for a role
type RolePermissions struct {
	RoleID string
	Perms  []Permission
}

// RBAC structure holds the core data for role-based-access control.
type RBAC struct {
	mu                   sync.RWMutex
	userRolesCache       map[string]UserRoles
	rolePermissionsCache map[string]RolePermissions
	fetchUserRoles       func(userID string) (UserRoles, error)
	fetchRolePerms       func(roleID string) (RolePermissions, error)
}

// NewRBAC creates an RBAC instance with provided fetch functions
func NewRBAC(
	fetchUserRoles func(string) (UserRoles, error),
	fetchRolePermissions func(string) (RolePermissions, error),
) *RBAC {
	// log.Trace("Creating a ne winstance of RBAC...")
	return &RBAC{
		userRolesCache:       make(map[string]UserRoles),
		rolePermissionsCache: make(map[string]RolePermissions),
		fetchUserRoles:       fetchUserRoles,
		fetchRolePerms:       fetchRolePermissions,
	}
}

// GetUserRoles fetches roles for a user, either from cache or by calling the fetch function
func (r *RBAC) GetUserRoles(userID string) (UserRoles, error) {
	r.mu.RLock()
	userRoles, found := r.userRolesCache[userID]
	r.mu.RUnlock()
	if found {
		return userRoles, nil
	}

	// Not found in cache, fetch and cache
	newUserRoles, err := r.fetchUserRoles(userID)
	if err != nil {
		return UserRoles{}, err
	}

	r.mu.Lock()
	r.userRolesCache[userID] = newUserRoles
	r.mu.Unlock()

	return newUserRoles, nil
}

// GetRolePermissions fetches permissions for a role, either from cache or by calling the fetch function
func (r *RBAC) GetRolePermissions(roleID string) (RolePermissions, error) {
	r.mu.RLock()
	rolePerms, found := r.rolePermissionsCache[roleID]
	r.mu.RUnlock()
	if found {
		return rolePerms, nil
	}

	// Not found in cache, fetch and cache
	newRolePerms, err := r.fetchRolePerms(roleID)
	if err != nil {
		return RolePermissions{}, err
	}

	r.mu.Lock()
	r.rolePermissionsCache[roleID] = newRolePerms
	r.mu.Unlock()

	return newRolePerms, nil
}
