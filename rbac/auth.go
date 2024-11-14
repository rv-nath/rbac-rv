package rbac

// CheckPermission checks if a user has permission to perform an action on a specific resource instance within a given scope.
func (r *RBAC) CheckPermission(userID, resourceID, resourceName, action, scope string) bool {
	// Retrieve roles associated with the user
	userRoles, err := r.GetUserRoles(userID)
	if err != nil {
		return false // Log error or handle appropriately
	}

	for _, roleID := range userRoles.Roles {
		// Retrieve permissions associated with the role
		rolePerms, err := r.GetRolePermissions(roleID)
		if err != nil {
			continue // Log error or handle appropriately
		}

		// Iterate over each permission to find a match
		for _, perm := range rolePerms.Perms {
			if perm.Resource == resourceName && perm.Permission == action && perm.Allowed {
				// Additional scope-based checks
				switch scope {
				case "self":
					if r.isUserOwnerOfResource(userID, resourceID) {
						return true
					}
				case "group":
					if r.isUserInGroupWithResourceOwner(userID, resourceID) {
						return true
					}
				case "global":
					return true // Global scope grants permission directly
				}
			}
		}
	}
	return false
}

// isUserOwnerOfResource checks if a user owns a specific resource instance (pseudo-check).
func (r *RBAC) isUserOwnerOfResource(userID, resourceID string) bool {
	// Implement actual ownership check logic here, e.g., database query
	return true // Placeholder
}

// isUserInGroupWithResourceOwner checks if a user shares a group with the owner of a specific resource instance (pseudo-check).
func (r *RBAC) isUserInGroupWithResourceOwner(userID, resourceID string) bool {
	// Implement actual group membership check logic here, e.g., database query
	return true // Placeholder
}
