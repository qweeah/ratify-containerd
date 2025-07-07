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

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// RatifyOutput represents the structure of ratify verify command output
type RatifyOutput struct {
	IsSuccess bool        `json:"isSuccess"`
	Result    interface{} `json:"result,omitempty"`
	Error     string      `json:"error,omitempty"`
}

var (
	name           string
	digest         string
	stdinMediaType string
)

func init() {
	// Define command line flags
	flag.StringVar(&name, "name", "", "Container image name (required)")
	flag.StringVar(&digest, "digest", "", "Container image digest (required)")
	flag.StringVar(&stdinMediaType, "stdin-media-type", "", "Stdin media type")
}

func main() {
	// Parse command line flags
	flag.Parse()

	// Validate required flags
	if name == "" {
		fmt.Fprintf(os.Stderr, "Error: -name flag is required\n")
		flag.Usage()
		os.Exit(1)
	}
	if digest == "" {
		fmt.Fprintf(os.Stderr, "Error: -digest flag is required\n")
		flag.Usage()
		os.Exit(1)
	}

	// Set HOME environment variable
	os.Setenv("HOME", "/root")

	// Set default paths
	configPath := filepath.Join("/root", ".ratify", "config.json")

	// Check if ratify config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		fmt.Println("ratify config file not found. failing open")
		os.Exit(0)
	}

	// Execute ratify verify command
	ratifyOutput, err := executeRatifyVerify()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to execute ratify verify: %v\n", err)
		os.Exit(1)
	}

	// Parse and display the output
	var output RatifyOutput
	if err := json.Unmarshal([]byte(ratifyOutput), &output); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to parse ratify output: %v\n", err)
		fmt.Println("Raw output:", ratifyOutput)
		os.Exit(1)
	}

	// Pretty print the JSON output
	prettyJSON, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		fmt.Println("Raw output:", ratifyOutput)
	} else {
		fmt.Println(string(prettyJSON))
	}

	// Check if verification succeeded
	if output.IsSuccess {
		fmt.Println("ratify verification succeeded")
		os.Exit(0)
	} else {
		fmt.Println("ratify verification failed")
		os.Exit(1)
	}
}

func executeRatifyVerify() (string, error) {
	// Set default paths
	configPath := filepath.Join("/root", ".ratify", "config.json")
	ratifyBin := filepath.Join("/root", ".ratify", "bin", "ratify")

	// Construct the ratify verify command arguments
	args := []string{"verify", "-c", configPath, "-s", name, "--digest", digest}

	// Add stdin-media-type if provided
	if stdinMediaType != "" {
		args = append(args, "--stdin-media-type", stdinMediaType)
	}

	// Create the command
	cmd := exec.Command(ratifyBin, args...)

	// Set environment variables
	cmd.Env = append(os.Environ(), "HOME=/root")

	// Execute the command and capture output
	output, err := cmd.Output()
	if err != nil {
		// ratify might exit with non-zero code even for valid verification failures
		// so we try to get the output from stderr as well
		if exitError, ok := err.(*exec.ExitError); ok {
			// Combine stdout and stderr
			combined := string(output) + string(exitError.Stderr)
			if len(combined) > 0 {
				return combined, nil
			}
		}
		return "", fmt.Errorf("ratify command failed: %v", err)
	}

	return string(output), nil
}
