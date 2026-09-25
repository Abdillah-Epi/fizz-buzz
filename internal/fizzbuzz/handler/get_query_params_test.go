package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	fizzbuzzDto "github.com/Abdillah-Epi/fizz-buzz/internal/fizzbuzz/dto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetQueryParams(t *testing.T) {
	tests := []struct {
		name    string
		query   string
		want    fizzbuzzDto.FizzBuzzRequest
		wantErr string
	}{
		{
			name:  "parses every parameter",
			query: "int1=3&int2=5&limit=15&str1=Fizz&str2=Buzz",
			want:  fizzbuzzDto.FizzBuzzRequest{Int1: 3, Int2: 5, Limit: 15, Str1: "Fizz", Str2: "Buzz"},
		},
		{
			name:  "parameter order does not matter",
			query: "str2=Buzz&limit=15&str1=Fizz&int2=5&int1=3",
			want:  fizzbuzzDto.FizzBuzzRequest{Int1: 3, Int2: 5, Limit: 15, Str1: "Fizz", Str2: "Buzz"},
		},
		{
			name:  "extra parameters are ignored",
			query: "int1=3&int2=5&limit=15&str1=Fizz&str2=Buzz&unknown=1",
			want:  fizzbuzzDto.FizzBuzzRequest{Int1: 3, Int2: 5, Limit: 15, Str1: "Fizz", Str2: "Buzz"},
		},
		{
			name:  "url encoded strings are decoded",
			query: "int1=3&int2=5&limit=15&str1=foo%20bar&str2=%C3%A9",
			want:  fizzbuzzDto.FizzBuzzRequest{Int1: 3, Int2: 5, Limit: 15, Str1: "foo bar", Str2: "é"},
		},
		{
			// Range checks live in the service, the parser only types the input.
			name:  "negative integers are parsed as-is",
			query: "int1=-3&int2=5&limit=15&str1=Fizz&str2=Buzz",
			want:  fizzbuzzDto.FizzBuzzRequest{Int1: -3, Int2: 5, Limit: 15, Str1: "Fizz", Str2: "Buzz"},
		},
		{
			name:  "empty strings are parsed as empty",
			query: "int1=3&int2=5&limit=15&str1=&str2=",
			want:  fizzbuzzDto.FizzBuzzRequest{Int1: 3, Int2: 5, Limit: 15},
		},
		{
			name:    "missing int1",
			query:   "int2=5&limit=15&str1=Fizz&str2=Buzz",
			wantErr: "invalid int1",
		},
		{
			name:    "non numeric int1",
			query:   "int1=three&int2=5&limit=15&str1=Fizz&str2=Buzz",
			wantErr: "invalid int1",
		},
		{
			name:    "fractional int1",
			query:   "int1=3.5&int2=5&limit=15&str1=Fizz&str2=Buzz",
			wantErr: "invalid int1",
		},
		{
			name:    "missing int2",
			query:   "int1=3&limit=15&str1=Fizz&str2=Buzz",
			wantErr: "invalid int2",
		},
		{
			name:    "non numeric int2",
			query:   "int1=3&int2=abc&limit=15&str1=Fizz&str2=Buzz",
			wantErr: "invalid int2",
		},
		{
			name:    "missing limit",
			query:   "int1=3&int2=5&str1=Fizz&str2=Buzz",
			wantErr: "invalid limit",
		},
		{
			name:    "non numeric limit",
			query:   "int1=3&int2=5&limit=lots&str1=Fizz&str2=Buzz",
			wantErr: "invalid limit",
		},
		{
			name:    "no parameters at all",
			query:   "",
			wantErr: "invalid int1",
		},
		{
			name:    "duplicate parameters use the first value",
			query:   "int1=3&int1=notanumber&int2=5&limit=15&str1=Fizz&str2=Buzz",
			want:    fizzbuzzDto.FizzBuzzRequest{Int1: 3, Int2: 5, Limit: 15, Str1: "Fizz", Str2: "Buzz"},
			wantErr: "",
		},
	}

	require := require.New(t)
	assert := assert.New(t)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			request := httptest.NewRequest(http.MethodGet, "/api/v1/fizzbuzz?"+tt.query, nil)

			got, err := getQueryParams(request)

			if tt.wantErr != "" {
				require.Error(err)
				assert.EqualError(err, tt.wantErr)
				assert.Equal(fizzbuzzDto.FizzBuzzRequest{}, got)
				return
			}

			require.NoError(err)
			assert.Equal(tt.want, got)
		})
	}
}
