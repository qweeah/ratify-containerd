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

package models

import (
	"encoding/json"
	"fmt"
	"os"
)

// Example usage of the optimized scoped configuration
// This demonstrates how another process would efficiently check scopes

// LoadScopedConfig loads the optimized scoped configuration from file
func LoadScopedConfig(filePath string) (*ScopedConfigOptimized, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}

	var config ScopedConfigOptimized
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config JSON: %v", err)
	}

	return &config, nil
}

// CheckScope checks if a scope exists in the configuration
// This is an O(1) operation thanks to the scopeMap
func CheckScope(config *ScopedConfigOptimized, scope string) bool {
	return config.HasScope(scope)
}

// ExampleUsage demonstrates how to use the optimized configuration
func ExampleUsage() {
	// Load the configuration
	config, err := LoadScopedConfig("/shared-data/ratify-config.json")
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		return
	}

	// Example scopes to check
	testScopes := []string{
		"mcr.microsoft.com/azurearck8s/metrics-agent",
		"docker.io/library/nginx",
		"mcr.microsoft.com/some/other/image",
	}

	// Check each scope efficiently (O(1) per check)
	for _, scope := range testScopes {
		if CheckScope(config, scope) {
			fmt.Printf("✓ Scope '%s' is verified\n", scope)
		} else {
			fmt.Printf("✗ Scope '%s' is not in verification list\n", scope)
		}
	}

	fmt.Printf("Configuration contains %d total scopes\n", len(config.ScopeMap))
	fmt.Printf("Last updated: %s\n", config.LastUpdated.Format("2006-01-02 15:04:05"))
}
