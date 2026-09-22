package config

import (
	"fmt"
	"os"
	"strconv"
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

func Load() (Config, error) {
	maxLimit, err := getIntEnv("FIZZBUZZ_MAX_LIMIT", 10000)
	if err != nil {
		return Config{}, fmt.Errorf("load FIZZBUZZ_MAX_LIMIT: %w", err)
	}

	if maxLimit <= 0 {
		return Config{}, fmt.Errorf("FIZZBUZZ_MAX_LIMIT must be greater than 0")
	}

	clickHouse, err := loadClickHouseConfig()
	if err != nil {
		return Config{}, fmt.Errorf("load clickhouse config: %w", err)
	}

	return Config{
		FizzBuzz: FizzBuzzConfig{
			MaxLimit: maxLimit,
		},
		ClickHouse: clickHouse,
	}, nil
}

func loadClickHouseConfig() (ClickHouseConfig, error) {
	port, err := getIntEnv("CLICKHOUSE_PORT", 9000)
	if err != nil {
		return ClickHouseConfig{}, fmt.Errorf("load CLICKHOUSE_PORT: %w", err)
	}

	host := getEnv("CLICKHOUSE_HOST", "localhost")
	database := getEnv("CLICKHOUSE_DB", "fizzbuzz")
	username := getEnv("CLICKHOUSE_USER", "fizzbuzz")
	password := getEnv("CLICKHOUSE_PASSWORD", "")

	if host == "" {
		return ClickHouseConfig{}, fmt.Errorf("CLICKHOUSE_HOST must not be empty")
	}
	if database == "" {
		return ClickHouseConfig{}, fmt.Errorf("CLICKHOUSE_DB must not be empty")
	}
	if username == "" {
		return ClickHouseConfig{}, fmt.Errorf("CLICKHOUSE_USER must not be empty")
	}

	return ClickHouseConfig{
		Host:     host,
		Port:     port,
		Database: database,
		Username: username,
		Password: password,
	}, nil
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)

	if value == "" {
		return fallback
	}
	return value
}

func getIntEnv(key string, fallback int) (int, error) {
	value := os.Getenv(key)

	if value == "" {
		return fallback, nil
	}

	result, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("%s must be a valid integer: %w", key, err)
	}

	return result, nil
}
