package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var requiredKeys = []string{
	"CLICKHOUSE_HOST",
	"CLICKHOUSE_DB",
	"CLICKHOUSE_USER",
	"CLICKHOUSE_PASSWORD",
}

var optionalKeys = []string{
	"FIZZBUZZ_MAX_LIMIT",
	"CLICKHOUSE_PORT",
}

func unsetEnv(t *testing.T, key string) {
	require := require.New(t)
	t.Helper()

	original, existed := os.LookupEnv(key)
	require.NoError(os.Unsetenv(key))

	t.Cleanup(func() {
		if existed {
			_ = os.Setenv(key, original)
			return
		}
		_ = os.Unsetenv(key)
	})
}

func setValidEnv(t *testing.T) {
	t.Helper()

	unsetEnv(t, "FIZZBUZZ_MAX_LIMIT")
	unsetEnv(t, "CLICKHOUSE_PORT")

	t.Setenv("CLICKHOUSE_HOST", "localhost")
	t.Setenv("CLICKHOUSE_DB", "fizzbuzz")
	t.Setenv("CLICKHOUSE_USER", "fizzbuzz")
	t.Setenv("CLICKHOUSE_PASSWORD", "fizzbuzz")
}

func TestLoadUsesDefaultsForOptionalValues(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)
	setValidEnv(t)

	cfg, err := Load()

	require.NoError(err)
	assert.Equal(10000, cfg.FizzBuzz.MaxLimit)
	assert.Equal(ClickHouseConfig{
		Host:     "localhost",
		Port:     9000,
		Database: "fizzbuzz",
		Username: "fizzbuzz",
		Password: "fizzbuzz",
	}, cfg.ClickHouse)
}

func TestLoadReadsEveryValueFromTheEnvironment(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)

	setValidEnv(t)
	t.Setenv("FIZZBUZZ_MAX_LIMIT", "250")
	t.Setenv("CLICKHOUSE_HOST", "clickhouse.internal")
	t.Setenv("CLICKHOUSE_PORT", "9100")
	t.Setenv("CLICKHOUSE_DB", "analytics")
	t.Setenv("CLICKHOUSE_USER", "reader")
	t.Setenv("CLICKHOUSE_PASSWORD", "secret")

	cfg, err := Load()

	require.NoError(err)
	assert.Equal(250, cfg.FizzBuzz.MaxLimit)
	assert.Equal(ClickHouseConfig{
		Host:     "clickhouse.internal",
		Port:     9100,
		Database: "analytics",
		Username: "reader",
		Password: "secret",
	}, cfg.ClickHouse)
}

func TestLoadRejectsBlankRequiredValues(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)

	for _, key := range requiredKeys {
		t.Run(key, func(t *testing.T) {
			setValidEnv(t)
			t.Setenv(key, "")
			_, err := Load()
			require.Error(err)
			assert.ErrorContains(err, key+" is required and must not be empty")
		})
	}
}

func TestLoadRejectsMissingRequiredValues(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)

	for _, key := range requiredKeys {
		t.Run(key, func(t *testing.T) {
			setValidEnv(t)
			unsetEnv(t, key)
			_, err := Load()
			require.Error(err)
			assert.ErrorContains(err, key+" is required and must not be empty")
		})
	}
}

func TestLoadRejectsWhitespaceOnlyRequiredValues(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)

	setValidEnv(t)
	t.Setenv("CLICKHOUSE_DB", "   ")
	_, err := Load()
	require.Error(err)
	assert.ErrorContains(err, "CLICKHOUSE_DB is required and must not be empty")
}

func TestLoadReportsEveryProblemAtOnce(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)

	for _, key := range append(requiredKeys, optionalKeys...) {
		unsetEnv(t, key)
	}
	_, err := Load()
	require.Error(err)
	for _, key := range requiredKeys {
		assert.ErrorContains(err, key, "every missing value should be reported together")
	}
}

func TestLoadRejectsInvalidValues(t *testing.T) {
	tests := []struct {
		name    string
		env     map[string]string
		wantErr []string
	}{
		{
			name:    "non numeric max limit",
			env:     map[string]string{"FIZZBUZZ_MAX_LIMIT": "many"},
			wantErr: []string{"FIZZBUZZ_MAX_LIMIT must be a valid integer"},
		},
		{
			name:    "zero max limit",
			env:     map[string]string{"FIZZBUZZ_MAX_LIMIT": "0"},
			wantErr: []string{"FIZZBUZZ_MAX_LIMIT must be greater than 0"},
		},
		{
			name:    "negative max limit",
			env:     map[string]string{"FIZZBUZZ_MAX_LIMIT": "-100"},
			wantErr: []string{"FIZZBUZZ_MAX_LIMIT must be greater than 0"},
		},
		{
			name:    "non numeric clickhouse port",
			env:     map[string]string{"CLICKHOUSE_PORT": "nine-thousand"},
			wantErr: []string{"CLICKHOUSE_PORT must be a valid integer"},
		},
		{
			name:    "zero clickhouse port",
			env:     map[string]string{"CLICKHOUSE_PORT": "0"},
			wantErr: []string{"CLICKHOUSE_PORT must be between 1 and 65535"},
		},
		{
			name:    "clickhouse port above the valid range",
			env:     map[string]string{"CLICKHOUSE_PORT": "70000"},
			wantErr: []string{"CLICKHOUSE_PORT must be between 1 and 65535"},
		},
		{
			name: "two problems are reported together",
			env: map[string]string{
				"FIZZBUZZ_MAX_LIMIT": "0",
				"CLICKHOUSE_PORT":    "70000",
			},
			wantErr: []string{
				"FIZZBUZZ_MAX_LIMIT must be greater than 0",
				"CLICKHOUSE_PORT must be between 1 and 65535",
			},
		},
	}

	require := require.New(t)
	assert := assert.New(t)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setValidEnv(t)
			for key, value := range tt.env {
				t.Setenv(key, value)
			}
			_, err := Load()
			require.Error(err)
			assert.ErrorContains(err, "invalid configuration")
			for _, want := range tt.wantErr {
				assert.ErrorContains(err, want)
			}
		})
	}
}

func TestRequiredEnv(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)

	unsetEnv(t, "FIZZBUZZ_TEST_ABSENT")

	_, err := requiredEnv("FIZZBUZZ_TEST_ABSENT")
	require.Error(err)
	assert.EqualError(err, "FIZZBUZZ_TEST_ABSENT is required and must not be empty")

	t.Setenv("FIZZBUZZ_TEST_BLANK", "")
	_, err = requiredEnv("FIZZBUZZ_TEST_BLANK")
	require.Error(err)

	t.Setenv("FIZZBUZZ_TEST_SET", "value")
	value, err := requiredEnv("FIZZBUZZ_TEST_SET")
	require.NoError(err)
	assert.Equal("value", value)
}

func TestOptionalIntEnv(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)

	unsetEnv(t, "FIZZBUZZ_TEST_ABSENT")
	t.Setenv("FIZZBUZZ_TEST_BLANK", "")
	t.Setenv("FIZZBUZZ_TEST_INT", "42")
	t.Setenv("FIZZBUZZ_TEST_ZERO", "0")
	t.Setenv("FIZZBUZZ_TEST_BAD", "4o2")

	value, err := optionalIntEnv("FIZZBUZZ_TEST_INT", 1)
	require.NoError(err)
	assert.Equal(42, value)

	value, err = optionalIntEnv("FIZZBUZZ_TEST_ABSENT", 7)
	require.NoError(err)
	assert.Equal(7, value, "an absent value falls back to the default")

	value, err = optionalIntEnv("FIZZBUZZ_TEST_BLANK", 7)
	require.NoError(err)
	assert.Equal(7, value, "a blank value falls back to the default")

	value, err = optionalIntEnv("FIZZBUZZ_TEST_ZERO", 1)
	require.NoError(err)
	assert.Equal(0, value)

	_, err = optionalIntEnv("FIZZBUZZ_TEST_BAD", 1)
	require.Error(err)
	assert.ErrorContains(err, "FIZZBUZZ_TEST_BAD must be a valid integer")
}
