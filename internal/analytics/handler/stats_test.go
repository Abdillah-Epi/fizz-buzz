package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	analyticsDto "github.com/Abdillah-Epi/fizz-buzz/internal/analytics/dto"
	"github.com/Abdillah-Epi/fizz-buzz/internal/analytics/model"
	"github.com/Abdillah-Epi/fizz-buzz/internal/analytics/repository"
	"github.com/Abdillah-Epi/fizz-buzz/internal/analytics/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeRepository struct {
	stats    model.Stats
	statsErr error
}

func (f *fakeRepository) RecordEvent(context.Context, model.RequestEvent) error {
	return nil
}

func (f *fakeRepository) GetStats(context.Context) (model.Stats, error) {
	if f.statsErr != nil {
		return model.Stats{}, f.statsErr
	}

	return f.stats, nil
}

var _ repository.Repository = (*fakeRepository)(nil)

func firstSeen() time.Time {
	return time.Date(2026, time.September, 23, 8, 33, 28, 0, time.UTC)
}

func lastSeen() time.Time {
	return time.Date(2026, time.September, 23, 8, 33, 36, 0, time.UTC)
}

func getStats(t *testing.T, repo *fakeRepository) *httptest.ResponseRecorder {
	t.Helper()

	recorder := httptest.NewRecorder()
	NewStatsHandler(service.New(repo)).GetStats(recorder, httptest.NewRequest(http.MethodGet, "/api/v1/stats", nil))

	return recorder
}

func TestGetStatsReturnsAggregatedStatistics(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)
	t.Parallel()

	repo := &fakeRepository{
		stats: model.Stats{
			TotalRequests:        1255,
			RequestsLastHour:     12,
			RequestsLast24Hours:  900,
			UniqueConfigurations: 6,
			MostUsed: model.MostUsedRequest{
				Int1:        15,
				Int2:        25,
				LimitValue:  350,
				Str1:        "Fezz",
				Str2:        "Byzz",
				Hits:        450,
				Share:       35.85657370517929,
				FirstSeenAt: firstSeen(),
				LastSeenAt:  lastSeen(),
			},
		},
	}

	recorder := getStats(t, repo)

	require.Equal(http.StatusOK, recorder.Code)
	assert.Equal("application/json", recorder.Header().Get("Content-Type"))

	var body analyticsDto.StatsResponse
	require.NoError(json.NewDecoder(recorder.Body).Decode(&body))

	assert.Equal(analyticsDto.StatsResponse{
		TotalRequests:        1255,
		RequestsLastHour:     12,
		RequestsLast24Hours:  900,
		UniqueConfigurations: 6,
		MostUsed: analyticsDto.MostUsedRequest{
			Int1:        15,
			Int2:        25,
			Limit:       350,
			Str1:        "Fezz",
			Str2:        "Byzz",
			Hits:        450,
			Share:       35.85657370517929,
			FirstSeenAt: firstSeen(),
			LastSeenAt:  lastSeen(),
		},
	}, body)
}

func TestGetStatsResponseUsesDocumentedJSONKeys(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)
	t.Parallel()

	recorder := getStats(t, &fakeRepository{})

	var raw map[string]json.RawMessage
	require.NoError(json.NewDecoder(recorder.Body).Decode(&raw))

	assert.ElementsMatch([]string{
		"total_requests",
		"requests_last_hour",
		"requests_last_24h",
		"unique_configurations",
		"most_used",
	}, keysOf(raw))

	var mostUsed map[string]json.RawMessage
	require.NoError(json.Unmarshal(raw["most_used"], &mostUsed))

	assert.ElementsMatch([]string{
		"int1",
		"int2",
		"limit",
		"str1",
		"str2",
		"hits",
		"share",
		"first_seen_at",
		"last_seen_at",
	}, keysOf(mostUsed))
}

func TestGetStatsWithoutEventsReturnsZeroes(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)
	t.Parallel()

	recorder := getStats(t, &fakeRepository{})

	require.Equal(http.StatusOK, recorder.Code)

	var body analyticsDto.StatsResponse
	require.NoError(json.NewDecoder(recorder.Body).Decode(&body))

	assert.Zero(body.TotalRequests)
	assert.Zero(body.RequestsLastHour)
	assert.Zero(body.RequestsLast24Hours)
	assert.Zero(body.UniqueConfigurations)
	assert.Zero(body.MostUsed.Hits)
	assert.Zero(body.MostUsed.Share)
}

func TestGetStatsReturnsInternalErrorOnRepositoryFailure(t *testing.T) {
	assert := assert.New(t)
	t.Parallel()

	recorder := getStats(t, &fakeRepository{statsErr: errors.New("clickhouse unavailable")})

	assert.Equal(http.StatusInternalServerError, recorder.Code)
	assert.Equal("application/json", recorder.Header().Get("Content-Type"))

	var body struct {
		Error string `json:"error"`
	}
	require.NoError(t, json.NewDecoder(recorder.Body).Decode(&body))
	assert.Equal("failed to get statistics", body.Error)
}

func keysOf(m map[string]json.RawMessage) []string {
	keys := make([]string, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}

	return keys
}
