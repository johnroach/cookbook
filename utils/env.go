package utils

import "os"

// GetEnv tries to find the environment key if it can't find it will fallback to a value defined.
func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}