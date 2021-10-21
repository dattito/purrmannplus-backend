package utils

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Returns the value of the environment variable named by the key, returns an error if the variable is not set
func GetEnvElseError(key string) (string, error) {
	if value, exists := os.LookupEnv(key); exists {
		return value, nil
	}
	return "", fmt.Errorf("environment variable not set: %v", key)
}

// Returns the value of the environment variable named by the key, returns the default value if the variable is not set
func GetBoolEnvElseError(key string) (bool, error) {
	if _, exists := os.LookupEnv(key); exists {
		return GetBoolEnv(key, false)
	}
	return false, fmt.Errorf("environment variable not set: %v", key)
}

// Returns the value of the environment variable named by the key, returns the default value if the variable is not set
func GetEnv(key, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}

// Returns the value of the environment variable named by the key, returns the default value if the variable is not set
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

// Returns the value of the environment variable named by the key, returns the default value if the variable is not set
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

// Development mode: Returns the value of the environment variable named by the key, returns the default value if the variable is not set
// Production mode: Returns the value of the environment variable named by the key, returns an error if the variable is not set
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

// Development mode: Returns the value of the environment variable named by the key, returns the default value if the variable is not set
// Production mode: Returns the value of the environment variable named by the key, returns an error if the variable is not set
func GetBoolEnvInDev(key string, defaultVal bool) (bool, error) {
	b, err := GetBoolEnv("PRODUCTION", false)

	if err != nil {
		return false, err
	}

	if b {
		return GetBoolEnvElseError(key)
	}

	return GetBoolEnv(key, defaultVal)
}

// Loads all environment variables from the .env file
func LoadDotEnvFile() error {
	return godotenv.Load(".env")
}

// Development: Gets the value of the environment variable named by the key, returns the default development value if the variable is not set
// Production: Gets the value of the environment variable named by the key, returns the default production value if the variable is not set
func GetBoolDevProdEnv(key string, devDefaultVal, prodDefaultVal bool) (bool, error) {
	productionMode, err := GetBoolEnv("PRODUCTION", false)

	if err != nil {
		return false, err
	}

	if productionMode {
		return GetBoolEnv(key, prodDefaultVal)
	}

	return GetBoolEnv(key, devDefaultVal)
}
