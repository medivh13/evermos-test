package env

import (
	"os"
)

func GetEnvironment() string {
	envList := map[string]string{
		"production":  "production",
		"staging":     "staging",
		"alpha":       "alpha",
		"development": "development",
	}

	env := os.Getenv("SRVENV")
	if val, ok := envList[env]; ok {
		return val
	}

	return "development"
}

func IsProduction() bool {
	return GetEnvironment() == "production"
}

func IsStaging() bool {
	return GetEnvironment() == "staging"
}

func IsDevelopment() bool {
	return GetEnvironment() == "development"
}

func ValidateEnvironment(onDevelopment, onStaging, onProduction bool) bool {
	return (IsDevelopment() && onDevelopment) || (IsStaging() && onStaging) || (IsProduction() && onProduction)
}
