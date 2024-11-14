package rbac

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

// Global cache for resources, initialized as empty
var (
	resourceCache     []string
	resourceCacheLock sync.RWMutex // Lock to protect concurrent access
)

// FetchResourcesCallback defines a function type for fetching resource names
type FetchResourcesCallback func() ([]string, error)

// UUID regex pattern
// var uuidRegex = regexp.MustCompile(`^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}$`)
var uuidRegex = regexp.MustCompile(`^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}$`)

// DetermineIntent infers the action and resource from the HTTP method and path.
// It returns the action, resource name, resource ID (if applicable), and an error if something goes wrong.
func DetermineIntent(method, path string, fetchResources FetchResourcesCallback) (string, string, string, error) {
	// Map HTTP methods to actions
	var action string
	switch method {
	case "GET":
		action = "view"
	case "POST":
		action = "create"
	case "PUT":
		action = "edit"
	case "DELETE":
		action = "delete"
	default:
		return "", "", "", errors.New("unsupported HTTP method")
	}

	// Populate resource cache if it's empty
	resourceCacheLock.RLock()
	isCacheEmpty := len(resourceCache) == 0
	resourceCacheLock.RUnlock()

	if isCacheEmpty {
		resources, err := fetchResources()
		if err != nil {
			return "", "", "", errors.New("failed to fetch resources")
		}
		resourceCacheLock.Lock()
		resourceCache = resources
		resourceCacheLock.Unlock()
	}

	// Split path into segments
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) < 1 {
		return "", "", "", errors.New("invalid path format")
	}

	// Iterate over URL segments to identify resource and resource ID
	resource, resourceID := "", ""
	for i, segment := range parts {
		// Check if the segment (singularized) matches any resource in the cache
		singularSegment := strings.TrimSuffix(segment, "s")

		resourceCacheLock.RLock()
		isResource := contains(resourceCache, singularSegment)
		resourceCacheLock.RUnlock()

		if isResource {
			resource = singularSegment
			// If there's another segment, treat it as resource ID, but validate
			if i+1 < len(parts) {
				nextSegment := parts[i+1]
				if isNumeric(nextSegment) || isUUID(nextSegment) {
					resourceID = nextSegment
				} else {
					return "", "", "", errors.New("invalid resource ID format; expected numeric or UUID")
				}
			}
			break
		}
	}

	if resource == "" {
		return "", "", "", errors.New("resource not found in path")
	}

	return action, resource, resourceID, nil
}

// Helper function to check if a slice contains a given string
func contains(slice []string, item string) bool {
	for _, elem := range slice {
		if elem == item {
			return true
		}
	}
	return false
}

// Helper function to check if a string is numeric
func isNumeric(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}

// Helper function to check if a string matches the UUID format
func isUUID(s string) bool {
	return uuidRegex.MatchString(s)
}
