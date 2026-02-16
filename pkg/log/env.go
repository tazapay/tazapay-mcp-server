//nolint:forbidigo // object url will be a variable url
package log

import "os"

// Get returns the value of the environment variable named by the key.
func Get(key string) string {
	return os.Getenv(key)
}

// Set sets the value of the environment variable named by the key.
func Set(key, value string) error {
	return os.Setenv(key, value)
}
