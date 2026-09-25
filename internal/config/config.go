package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	defaultMaxLimit       = 10000
	defaultClickHousePort = 9000

	minPort = 1
	maxPort = 65535
)

type Config struct {
	FizzBuzz   FizzBuzzConfig
	ClickHouse ClickHouseConfig
}

type FizzBuzzConfig struct {
	MaxLimit int
}

type ClickHouseConfig struct {
	Host     string
	Port     int
	Database string
	Username string
	Password string
}

// Load builds the configuration from the environment.
//
// Deployment-specific values (the ClickHouse connection identity) are required:
// when one is missing or blank it is reported instead of being silently replaced
// by a guessed default, so a misconfigured process fails at startup rather than
// connecting to the wrong database. Only values with a safe default are optional.
//
// Every problem is reported at once, so a deployment can be fixed in one pass.
func Load() (Config, error) {
	var problems []error

	maxLimit, err := optionalIntEnv("FIZZBUZZ_MAX_LIMIT", defaultMaxLimit)
	switch {
	case err != nil:
		problems = append(problems, err)
	case maxLimit <= 0:
		problems = append(problems, fmt.Errorf("FIZZBUZZ_MAX_LIMIT must be greater than 0, got %d", maxLimit))
	}

	clickHouse, err := loadClickHouseConfig()
	if err != nil {
		problems = append(problems, err)
	}

	if err := errors.Join(problems...); err != nil {
		return Config{}, fmt.Errorf("invalid configuration: %w", err)
	}

	return Config{
		FizzBuzz: FizzBuzzConfig{
			MaxLimit: maxLimit,
		},
		ClickHouse: clickHouse,
	}, nil
}

func loadClickHouseConfig() (ClickHouseConfig, error) {
	var problems []error

	host, err := requiredEnv("CLICKHOUSE_HOST")
	if err != nil {
		problems = append(problems, err)
	}

	database, err := requiredEnv("CLICKHOUSE_DB")
	if err != nil {
		problems = append(problems, err)
	}

	username, err := requiredEnv("CLICKHOUSE_USER")
	if err != nil {
		problems = append(problems, err)
	}

	password, err := requiredEnv("CLICKHOUSE_PASSWORD")
	if err != nil {
		problems = append(problems, err)
	}

	port, err := optionalIntEnv("CLICKHOUSE_PORT", defaultClickHousePort)
	switch {
	case err != nil:
		problems = append(problems, err)
	case port < minPort || port > maxPort:
		problems = append(problems, fmt.Errorf("CLICKHOUSE_PORT must be between %d and %d, got %d", minPort, maxPort, port))
	}

	if err := errors.Join(problems...); err != nil {
		return ClickHouseConfig{}, err
	}

	return ClickHouseConfig{
		Host:     host,
		Port:     port,
		Database: database,
		Username: username,
		Password: password,
	}, nil
}

func requiredEnv(key string) (string, error) {
	value, ok := os.LookupEnv(key)
	if !ok || strings.TrimSpace(value) == "" {
		return "", fmt.Errorf("%s is required and must not be empty", key)
	}

	return value, nil
}

func optionalIntEnv(key string, fallback int) (int, error) {
	value, ok := os.LookupEnv(key)
	if !ok || strings.TrimSpace(value) == "" {
		return fallback, nil
	}

	parsed, err := strconv.Atoi(strings.TrimSpace(value))
	if err != nil {
		return 0, fmt.Errorf("%s must be a valid integer, got %q", key, value)
	}

	return parsed, nil
}
