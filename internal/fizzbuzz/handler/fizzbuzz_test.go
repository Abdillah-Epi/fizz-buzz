package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	analyticsModel "github.com/Abdillah-Epi/fizz-buzz/internal/analytics/model"
	analyticsService "github.com/Abdillah-Epi/fizz-buzz/internal/analytics/service"
	"github.com/Abdillah-Epi/fizz-buzz/internal/config"
	fizzbuzzDto "github.com/Abdillah-Epi/fizz-buzz/internal/fizzbuzz/dto"
	fizzbuzzService "github.com/Abdillah-Epi/fizz-buzz/internal/fizzbuzz/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testMaxLimit = 100

const validQuery = "int1=3&int2=5&limit=15&str1=Fizz&str2=Buzz"

type fakeAnalyticsRepository struct {
	events    []analyticsModel.RequestEvent
	recordErr error
}

func (f *fakeAnalyticsRepository) RecordEvent(_ context.Context, event analyticsModel.RequestEvent) error {
	if f.recordErr != nil {
		return f.recordErr
	}

	f.events = append(f.events, event)

	return nil
}

func (f *fakeAnalyticsRepository) GetStats(context.Context) (analyticsModel.Stats, error) {
	return analyticsModel.Stats{}, nil
}

func newTestHandler(repo *fakeAnalyticsRepository) *FizzBuzzHandler {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	service := fizzbuzzService.NewFizzBuzzService(config.Config{
		FizzBuzz: config.FizzBuzzConfig{MaxLimit: testMaxLimit},
	})

	return NewFizzBuzzHandler(service, analyticsService.New(repo), logger)
}

func generate(handler *FizzBuzzHandler, query string) *httptest.ResponseRecorder {
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api/v1/fizzbuzz?"+query, nil)

	handler.Generate(recorder, request)

	return recorder
}

func TestGenerateReturnsSequence(t *testing.T) {
	t.Parallel()
	require := require.New(t)
	assert := assert.New(t)

	recorder := generate(newTestHandler(&fakeAnalyticsRepository{}), validQuery)

	assert.Equal(http.StatusOK, recorder.Code)
	assert.Equal("application/json", recorder.Header().Get("Content-Type"))

	var body fizzbuzzDto.FizzBuzzResponse
	require.NoError(json.NewDecoder(recorder.Body).Decode(&body))

	expected := []string{"1", "2", "Fizz", "4", "Buzz", "Fizz", "7", "8", "Fizz", "Buzz", "11", "Fizz", "13", "14", "FizzBuzz"}
	assert.Equal(expected, body.Values)
}

func TestGenerateRecordsAnalyticsEvent(t *testing.T) {
	t.Parallel()
	require := require.New(t)
	assert := assert.New(t)

	repo := &fakeAnalyticsRepository{}

	before := time.Now().UTC()
	recorder := generate(newTestHandler(repo), validQuery)
	after := time.Now().UTC()

	require.Equal(http.StatusOK, recorder.Code)
	require.Len(repo.events, 1)

	event := repo.events[0]
	assert.Equal(uint32(3), event.Int1)
	assert.Equal(uint32(5), event.Int2)
	assert.Equal(uint32(15), event.LimitValue)
	assert.Equal("Fizz", event.Str1)
	assert.Equal("Buzz", event.Str2)
	assert.Equal(time.UTC, event.Timestamp.Location())
	assert.WithinDuration(before, event.Timestamp, time.Second)
	assert.LessOrEqual(event.Timestamp, after)
}

func TestGenerateRejectsInvalidInput(t *testing.T) {
	tests := []struct {
		name     string
		query    string
		wantBody string
	}{
		{
			name:     "missing int1",
			query:    "int2=5&limit=15&str1=Fizz&str2=Buzz",
			wantBody: "invalid int1",
		},
		{
			name:     "non numeric limit",
			query:    "int1=3&int2=5&limit=lots&str1=Fizz&str2=Buzz",
			wantBody: "invalid limit",
		},
		{
			name:     "zero int2 is rejected by the service",
			query:    "int1=3&int2=0&limit=15&str1=Fizz&str2=Buzz",
			wantBody: "int2 must be greater than 0",
		},
		{
			name:     "limit above the maximum is rejected by the service",
			query:    "int1=3&int2=5&limit=101&str1=Fizz&str2=Buzz",
			wantBody: "limit must be less than or equal to 100",
		},
		{
			name:     "divisor too large for the analytics column is rejected",
			query:    "int1=4294967296&int2=5&limit=15&str1=Fizz&str2=Buzz",
			wantBody: "int1 must be less than or equal to 4294967295",
		},
		{
			name:     "empty str1 is rejected by the service",
			query:    "int1=3&int2=5&limit=15&str1=&str2=Buzz",
			wantBody: "str1 must not be empty",
		},
		{
			name:     "empty str2 is rejected by the service",
			query:    "int1=3&int2=5&limit=15&str1=Fizz&str2=",
			wantBody: "str2 must not be empty",
		},
	}

	assert := assert.New(t)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := &fakeAnalyticsRepository{}

			recorder := generate(newTestHandler(repo), tt.query)

			assert.Equal(http.StatusBadRequest, recorder.Code)
			assert.Equal(tt.wantBody+"\n", recorder.Body.String())
			assert.Empty(repo.events, "rejected requests must not be recorded")
		})
	}
}

func TestGenerateRecordsLargeDivisorsWithoutTruncation(t *testing.T) {
	t.Parallel()
	require := require.New(t)
	assert := assert.New(t)

	repo := &fakeAnalyticsRepository{}
	query := fmt.Sprintf("int1=%d&int2=3&limit=6&str1=a&str2=b", fizzbuzzService.MaxDivisor)

	recorder := generate(newTestHandler(repo), query)

	require.Equal(http.StatusOK, recorder.Code)
	require.Len(repo.events, 1)
	assert.Equal(uint32(fizzbuzzService.MaxDivisor), repo.events[0].Int1)
	assert.Equal(uint32(3), repo.events[0].Int2)
}

func TestGenerateSucceedsWhenAnalyticsFails(t *testing.T) {
	t.Parallel()
	require := require.New(t)
	assert := assert.New(t)

	repo := &fakeAnalyticsRepository{recordErr: errors.New("clickhouse is down")}

	recorder := generate(newTestHandler(repo), validQuery)

	assert.Equal(t, http.StatusOK, recorder.Code)

	var body fizzbuzzDto.FizzBuzzResponse
	require.NoError(json.NewDecoder(recorder.Body).Decode(&body))
	assert.Len(body.Values, 15)
}

func TestGenerateReturnsInternalErrorWhenResponseCannotBeWritten(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	writer := &failingResponseWriter{header: http.Header{}}

	newTestHandler(&fakeAnalyticsRepository{}).Generate(writer, httptest.NewRequest(http.MethodGet, "/api/v1/fizzbuzz?"+validQuery, nil))

	assert.Equal(http.StatusInternalServerError, writer.status)
	assert.Contains(writer.body.String(), "failed to encode response")
}

type failingResponseWriter struct {
	header http.Header
	body   bytes.Buffer
	status int
}

func (w *failingResponseWriter) Header() http.Header {
	return w.header
}

func (w *failingResponseWriter) WriteHeader(statusCode int) {
	if w.status == 0 {
		w.status = statusCode
	}
}

func (w *failingResponseWriter) Write(p []byte) (int, error) {
	w.body.Write(p)

	return 0, errors.New("connection reset by peer")
}
