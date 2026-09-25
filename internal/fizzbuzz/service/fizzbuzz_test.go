package service

import (
	"fmt"
	"testing"

	"github.com/Abdillah-Epi/fizz-buzz/internal/config"
	"github.com/Abdillah-Epi/fizz-buzz/internal/fizzbuzz/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testMaxLimit = 100

func newTestService() *FizzBuzzService {
	return NewFizzBuzzService(config.Config{
		FizzBuzz: config.FizzBuzzConfig{MaxLimit: testMaxLimit},
	})
}

func TestGenerate(t *testing.T) {
	tests := []struct {
		name string
		req  model.Request
		want []string
	}{
		{
			name: "spec example: 3, 5 up to 15",
			req:  model.Request{Int1: 3, Int2: 5, Limit: 15, Str1: "Fizz", Str2: "Buzz"},
			want: []string{"1", "2", "Fizz", "4", "Buzz", "Fizz", "7", "8", "Fizz", "Buzz", "11", "Fizz", "13", "14", "FizzBuzz"},
		},
		{
			name: "limit of 1 yields a single number",
			req:  model.Request{Int1: 3, Int2: 5, Limit: 1, Str1: "Fizz", Str2: "Buzz"},
			want: []string{"1"},
		},
		{
			name: "int1 and int2 equal replace every multiple with the concatenation",
			req:  model.Request{Int1: 2, Int2: 2, Limit: 6, Str1: "a", Str2: "b"},
			want: []string{"1", "ab", "3", "ab", "5", "ab"},
		},
		{
			name: "identical strings still duplicate the token",
			req:  model.Request{Int1: 3, Int2: 3, Limit: 3, Str1: "x", Str2: "x"},
			want: []string{"1", "2", "xx"},
		},
		{
			name: "divisor of 1 replaces every value with its own string",
			req:  model.Request{Int1: 1, Int2: 5, Limit: 3, Str1: "A", Str2: "B"},
			want: []string{"A", "A", "A"},
		},
		{
			name: "two divisors of 1 concatenate on every value",
			req:  model.Request{Int1: 1, Int2: 1, Limit: 3, Str1: "A", Str2: "B"},
			want: []string{"AB", "AB", "AB"},
		},
		{
			name: "divisor greater than limit never matches",
			req:  model.Request{Int1: 50, Int2: 51, Limit: 3, Str1: "a", Str2: "b"},
			want: []string{"1", "2", "3"},
		},
		{
			name: "the largest storable divisor is accepted",
			req:  model.Request{Int1: int(MaxDivisor), Int2: 3, Limit: 6, Str1: "a", Str2: "b"},
			want: []string{"1", "2", "b", "4", "5", "b"},
		},
		{
			name: "multi-word and unicode strings are returned untouched",
			req:  model.Request{Int1: 2, Int2: 3, Limit: 6, Str1: "foo bar", Str2: "héllo"},
			want: []string{"1", "foo bar", "héllo", "foo bar", "5", "foo barhéllo"},
		},
	}

	require := require.New(t)
	assert := assert.New(t)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result, err := newTestService().Generate(tt.req)

			require.NoError(err)
			assert.Equal(tt.want, result.Values)
		})
	}
}

func TestGenerateAcceptsLimitAtMaxLimit(t *testing.T) {
	t.Parallel()
	require := require.New(t)
	assert := assert.New(t)

	req := model.Request{Int1: 3, Int2: 5, Limit: testMaxLimit, Str1: "Fizz", Str2: "Buzz"}

	result, err := newTestService().Generate(req)

	require.NoError(err)
	assert.Len(result.Values, testMaxLimit)
	assert.Equal("1", result.Values[0])
	assert.Equal("Fizz", result.Values[2])
	assert.Equal("FizzBuzz", result.Values[14])
	assert.Equal("Buzz", result.Values[testMaxLimit-1], "100 is a multiple of 5")
}

func TestGenerateValidationErrors(t *testing.T) {
	tests := []struct {
		name    string
		req     model.Request
		wantErr string
	}{
		{
			name:    "int1 must be strictly positive",
			req:     model.Request{Int1: 0, Int2: 5, Limit: 15, Str1: "Fizz", Str2: "Buzz"},
			wantErr: "int1 must be greater than 0",
		},
		{
			name:    "int1 must not be negative",
			req:     model.Request{Int1: -3, Int2: 5, Limit: 15, Str1: "Fizz", Str2: "Buzz"},
			wantErr: "int1 must be greater than 0",
		},
		{
			name:    "int2 must be strictly positive",
			req:     model.Request{Int1: 3, Int2: 0, Limit: 15, Str1: "Fizz", Str2: "Buzz"},
			wantErr: "int2 must be greater than 0",
		},
		{
			name:    "int2 must not be negative",
			req:     model.Request{Int1: 3, Int2: -5, Limit: 15, Str1: "Fizz", Str2: "Buzz"},
			wantErr: "int2 must be greater than 0",
		},
		{
			name:    "int1 must fit in the analytics column",
			req:     model.Request{Int1: int(MaxDivisor) + 1, Int2: 5, Limit: 15, Str1: "Fizz", Str2: "Buzz"},
			wantErr: fmt.Sprintf("int1 must be less than or equal to %d", MaxDivisor),
		},
		{
			name:    "int2 must fit in the analytics column",
			req:     model.Request{Int1: 3, Int2: int(MaxDivisor) + 1, Limit: 15, Str1: "Fizz", Str2: "Buzz"},
			wantErr: fmt.Sprintf("int2 must be less than or equal to %d", MaxDivisor),
		},
		{
			name:    "limit must be strictly positive",
			req:     model.Request{Int1: 3, Int2: 5, Limit: 0, Str1: "Fizz", Str2: "Buzz"},
			wantErr: "limit must be greater than 0",
		},
		{
			name:    "limit must not exceed the configured maximum",
			req:     model.Request{Int1: 3, Int2: 5, Limit: testMaxLimit + 1, Str1: "Fizz", Str2: "Buzz"},
			wantErr: fmt.Sprintf("limit must be less than or equal to %d", testMaxLimit),
		},
		{
			name:    "str1 must not be empty",
			req:     model.Request{Int1: 3, Int2: 5, Limit: 15, Str1: "", Str2: "Buzz"},
			wantErr: "str1 must not be empty",
		},
		{
			name:    "str2 must not be empty",
			req:     model.Request{Int1: 3, Int2: 5, Limit: 15, Str1: "Fizz", Str2: ""},
			wantErr: "str2 must not be empty",
		},
	}

	require := require.New(t)
	assert := assert.New(t)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result, err := newTestService().Generate(tt.req)

			require.Error(err)
			assert.EqualError(err, tt.wantErr)
			assert.Empty(result.Values, "no sequence should be produced when validation fails")
		})
	}
}

func TestGenerateReportsFirstFailingField(t *testing.T) {
	t.Parallel()
	require := require.New(t)
	assert := assert.New(t)

	req := model.Request{Int1: 0, Int2: 0, Limit: 0, Str1: "", Str2: ""}

	_, err := newTestService().Generate(req)
	require.Error(err)
	assert.EqualError(err, "int1 must be greater than 0")
}

func TestGenerateUsesConfiguredMaxLimit(t *testing.T) {
	t.Parallel()
	require := require.New(t)
	assert := assert.New(t)

	service := NewFizzBuzzService(config.Config{
		FizzBuzz: config.FizzBuzzConfig{MaxLimit: 10},
	})

	_, err := service.Generate(model.Request{Int1: 3, Int2: 5, Limit: 11, Str1: "Fizz", Str2: "Buzz"})
	require.Error(err)
	assert.EqualError(err, "limit must be less than or equal to 10")
}
