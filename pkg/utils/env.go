package utils

import (
	"os"
	"strings"
)

// GetEnv returns the value of the environment variable named key
func GetEnv(key string) string {
	return os.Getenv(key)
}

// IsProd returns true if the environment is production
// Example:
// ```go
// // export ENV=prod
//
//	if utils.IsProd() {
//		// do something
//	}
//
// ```
func IsProd() bool {
	for _, key := range []string{"ENV", "ENVIRONMENT"} {
		if v := strings.ToLower(os.Getenv(key)); v == "prod" || v == "production" {
			return true
		}
	}
	return false
}
