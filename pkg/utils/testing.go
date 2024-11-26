package utils

import (
	"os"
	"testing"

	"corvus_bot/pkg/config"
)

// IsTesting checks whether the application is running in test mode.
func IsTesting(cfg *config.Config) bool {
	// Check if the "GO_TEST" environment variable is set
	if os.Getenv("GO_TEST") == "1" {
		return true
	}

	// Check if the testing flag is set in the configuration
	if cfg != nil && cfg.Testing {
		return true
	}

	// Check if running under `go test`
	if testing.Short() {
		return true
	}

	return false
}
