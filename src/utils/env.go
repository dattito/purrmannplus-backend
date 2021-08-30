package utils

import (
	"fmt"
	"os"
	"strconv"
)

func GetEnvElseError(key string) (string, error) {
	if value, exists := os.LookupEnv(key); exists {
		return value, nil
	}
	return "", fmt.Errorf("environment variable not set: %v", key)
}

func GetEnv(key, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}

func GetIntEnv(key string, defaultVal int) (int, error) {
	if value, exists := os.LookupEnv(key); exists {
		intValue, err := strconv.Atoi(value)
		if err != nil {
			return 0, fmt.Errorf("can't convert enviroment variable to int: %v (Value: %v)", key, value)
		}
		return intValue, nil
	}
	return defaultVal, nil
}

func GetBoolEnv(key string, defaultVal bool) (bool, error) {
	if value, exists := os.LookupEnv(key); exists {
		boolValue, err := strconv.ParseBool(value)
		if err != nil {
			return defaultVal, fmt.Errorf("can't convert enviroment variable to bool: %v (Value: %v)", key, value)
		}
		return boolValue, nil
	}
	return defaultVal, nil
}

func GetEnvInDev(key, defaultVal string) (string, error) {
	b, err := GetBoolEnv("PRODUCTION", false)

	if err != nil {
		return "", err
	}

	if b {
		return GetEnvElseError(key)
	}

	return GetEnv(key, defaultVal), nil
}
