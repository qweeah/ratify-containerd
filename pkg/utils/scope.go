/*
Copyright The Ratify Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package utils

import (
	"sort"
)

// DeduplicateScopes removes duplicate scopes from a slice and returns a sorted slice
// This ensures consistent output and efficient processing
func DeduplicateScopes(scopes []string) []string {
	if len(scopes) == 0 {
		return scopes
	}
	
	// Use a map to track unique scopes
	scopeMap := make(map[string]struct{}, len(scopes))
	for _, scope := range scopes {
		if scope != "" { // skip empty strings
			scopeMap[scope] = struct{}{}
		}
	}
	
	// Convert back to slice
	deduplicated := make([]string, 0, len(scopeMap))
	for scope := range scopeMap {
		deduplicated = append(deduplicated, scope)
	}
	
	// Sort for consistent output
	sort.Strings(deduplicated)
	
	return deduplicated
}
