package utils

import "os"

func EnvIsTrue(key string) bool {
	val := os.Getenv(key)

	return val == "1" || val == "true"
}

func EnvIsRelease() bool {
	return os.Getenv("APP_MODE") == "release"
}
