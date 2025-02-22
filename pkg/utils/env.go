package utils

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/joho/godotenv"
)

// loadEnvOnce is a sync.Once instance
var loadEnvOnce sync.Once

// defaultEnvFiles is the list of env files to load in order of priority
var defaultEnvFiles = []string{".env.development", ".env.local", ".env", ".env.production", ".env.test"}

// Initenv loads the .env file
// Called Single time (Singleton)
func Initenv() error {
	var err error
	loadEnvOnce.Do(func() {
		if err = godotenv.Load(defaultEnvFiles...); err == nil {
			fmt.Printf("Loaded environment: %s\n", err)
			return
		}
	})
	return err
}

// Getenv returns the value of the environment variable named key
func Getenv(key string) string {
	if err := Initenv(); err != nil {
		fmt.Printf("Warning: Error loading environment: %v\n", err)
	}
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
