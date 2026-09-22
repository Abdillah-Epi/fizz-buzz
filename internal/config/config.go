package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	FizzBuzz FizzBuzzConfig
}

type FizzBuzzConfig struct {
	MaxLimit int
}

func Load() (Config, error) {
	maxLimit, err := getIntEnv("FIZZBUZZ_MAX_LIMIT", 10000)
	if err != nil {
		return Config{}, fmt.Errorf("load FIZZBUZZ_MAX_LIMIT: %w", err)
	}

	if maxLimit <= 0 {
		return Config{}, fmt.Errorf("FIZZBUZZ_MAX_LIMIT must be greater than 0")
	}

	return Config{
		FizzBuzz: FizzBuzzConfig{
			MaxLimit: maxLimit,
		},
	}, nil
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
