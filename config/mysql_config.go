package config

import (
	"os"
	"strconv"
)

// DatabaseConfig represents the MySQL database configuration.
type MySQLDatabaseConfig struct {
	Username     string
	Password     string
	Host         string
	Port         string
	Database     string
	MaxOpenConns int
	MaxIdleConns int
}

// GetConnectDatabseConfig returns the MySQL configuration read from environment variables.

func GetChoiceDBConfig() MySQLDatabaseConfig {
	return MySQLDatabaseConfig{
		Username:     getEnv("DB_USERNAME", "root"),
		Password:     getEnv("DB_PASSWORD", ""), // "password"
		Host:         getEnv("DB_HOST", "localhost"),
		Port:         getEnv("DB_PORT", "3306"),
		Database:     getEnv("DB_NAME", "shree_booking_ltd"), //"mydatabase"
		MaxOpenConns: getIntEnv("DB_MAX_OPEN_CONNS", 10),
		MaxIdleConns: getIntEnv("DB_MAX_IDLE_CONNS", 5),
	}
}

// getEnv is a helper function to get environment variables or return a default value.
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// getIntEnv is a helper function to get integer environment variables or return a default value.
func getIntEnv(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}
