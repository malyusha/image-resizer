package config

import (
	"os"
	"strings"
)

func IsTesting() bool {
	return strings.HasSuffix(os.Args[0], ".test")
}

func IsProd() bool {
	env := os.Getenv("ENV")

	return env == "" || env == "production"
}

func Is(env string) bool {
	return os.Getenv("ENV") == env
}
