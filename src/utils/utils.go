package utils

import (
	"log"
	"os"
	"strconv"
)

func GetEnvElseError(key string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	log.Fatalf("environment variable not set: %v", key)
	return ""
}

func GetEnv(key, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}

func GetIntEnv(key string, defaultVal int) int {
	if value, exists := os.LookupEnv(key); exists {
		intValue, err := strconv.Atoi(value)
		if err != nil {
			log.Fatalf("can't convert enviroment variable to int: %v (Value: %v)", key, value)
		}
		return intValue
	}
	return defaultVal
}

func GetBoolEnv(key string, defaultVal bool) bool {
	if value, exists := os.LookupEnv(key); exists {
		boolValue, err := strconv.ParseBool(value)
		if err != nil {
			log.Fatalf("can't convert enviroment variable to bool: %v (Value: %v)", key, value)
		}
		return boolValue
	}
	return defaultVal
}

func GetEnvInDev(key, defaultVal string) string {
	b := GetBoolEnv("PRODUCTION", false)

	if b {
		return GetEnvElseError(key)
	}

	return GetEnv(key, defaultVal)
}
