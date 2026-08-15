package config

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all service configuration parameters.
type Config struct {
	ServerPort          int
	DBHost              string
	DBPort              int
	DBUser              string
	DBPassword          string
	DBName              string
	DBSSLMode           string
	DBMaxConns          int32
	DBMinConns          int32
	DBMaxConnLifetime   time.Duration
	DBMaxConnIdleTime   time.Duration
	AggregationInterval time.Duration
}

// Load reads configuration from environment variables (loading .env if present) and applies defaults.
func Load() (*Config, error) {
	// Attempt to load .env file if it exists, ignore error if it does not
	_ = godotenv.Load()

	cfg := &Config{
		ServerPort:          getEnvAsInt("SERVER_PORT", 8080),
		DBHost:              getEnv("DB_HOST", "localhost"),
		DBPort:              getEnvAsInt("DB_PORT", 5432),
		DBUser:              getEnv("DB_USER", "postgres"),
		DBPassword:          getEnv("DB_PASSWORD", "postgres"),
		DBName:              getEnv("DB_NAME", "activity_tracker"),
		DBSSLMode:           getEnv("DB_SSLMODE", "disable"),
		DBMaxConns:          int32(getEnvAsInt("DB_MAX_CONNS", 25)),
		DBMinConns:          int32(getEnvAsInt("DB_MIN_CONNS", 5)),
		DBMaxConnLifetime:   getEnvAsDuration("DB_MAX_CONN_LIFETIME", 1*time.Hour),
		DBMaxConnIdleTime:   getEnvAsDuration("DB_MAX_CONN_IDLE_TIME", 30*time.Minute),
		AggregationInterval: getEnvAsDuration("AGGREGATION_INTERVAL", 4*time.Hour),
	}

	return cfg, nil
}

// DSN returns a standard PostgreSQL connection URL string.
func (c *Config) DSN() string {
	userInfo := url.UserPassword(c.DBUser, c.DBPassword)
	host := fmt.Sprintf("%s:%d", c.DBHost, c.DBPort)
	query := url.Values{}
	query.Set("sslmode", c.DBSSLMode)

	u := url.URL{
		Scheme:   "postgres",
		User:     userInfo,
		Host:     host,
		Path:     c.DBName,
		RawQuery: query.Encode(),
	}

	return u.String()
}

func getEnv(key, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok && val != "" {
		return val
	}
	return defaultVal
}

func getEnvAsInt(key string, defaultVal int) int {
	if val, ok := os.LookupEnv(key); ok && val != "" {
		if intVal, err := strconv.Atoi(val); err == nil {
			return intVal
		}
	}
	return defaultVal
}

func getEnvAsDuration(key string, defaultVal time.Duration) time.Duration {
	if val, ok := os.LookupEnv(key); ok && val != "" {
		if durationVal, err := time.ParseDuration(val); err == nil {
			return durationVal
		}
	}
	return defaultVal
}
