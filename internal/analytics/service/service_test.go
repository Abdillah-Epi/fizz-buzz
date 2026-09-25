package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Abdillah-Epi/fizz-buzz/internal/analytics/model"
	"github.com/Abdillah-Epi/fizz-buzz/internal/analytics/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeRepository struct {
	recordedEvents []model.RequestEvent
	recordErr      error
	stats          model.Stats
	statsErr       error
}

func (f *fakeRepository) RecordEvent(_ context.Context, event model.RequestEvent) error {
	if f.recordErr != nil {
		return f.recordErr
	}

	f.recordedEvents = append(f.recordedEvents, event)

	return nil
}

func (f *fakeRepository) GetStats(context.Context) (model.Stats, error) {
	if f.statsErr != nil {
		return model.Stats{}, f.statsErr
	}

	return f.stats, nil
}

var _ repository.Repository = (*fakeRepository)(nil)

func sampleEvent() model.RequestEvent {
	return model.RequestEvent{
		Timestamp:  time.Date(2026, time.September, 25, 10, 30, 0, 0, time.UTC),
		Int1:       3,
		Int2:       5,
		LimitValue: 100,
		Str1:       "fizz",
		Str2:       "buzz",
	}
}

func TestRecordEventDelegatesToRepository(t *testing.T) {
	t.Parallel()

	repo := &fakeRepository{}
	event := sampleEvent()

	require.NoError(t, New(repo).RecordEvent(context.Background(), event))

	require.Len(t, repo.recordedEvents, 1)
	assert.Equal(t, event, repo.recordedEvents[0])
}

func TestRecordEventPropagatesRepositoryError(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)
	t.Parallel()

	repo := &fakeRepository{recordErr: errors.New("insert failed")}

	err := New(repo).RecordEvent(context.Background(), sampleEvent())

	require.Error(err)
	assert.ErrorContains(err, "insert failed")
	assert.Empty(repo.recordedEvents)
}

func TestGetStatsReturnsRepositoryStats(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)
	t.Parallel()

	stats := model.Stats{
		TotalRequests:        42,
		RequestsLastHour:     10,
		RequestsLast24Hours:  40,
		UniqueConfigurations: 3,
		MostUsed: model.MostUsedRequest{
			Int1:       3,
			Int2:       5,
			LimitValue: 15,
			Str1:       "fizz",
			Str2:       "buzz",
			Hits:       20,
			Share:      47.6,
		},
	}
	repo := &fakeRepository{stats: stats}

	got, err := New(repo).GetStats(context.Background())

	require.NoError(err)
	assert.Equal(stats, got)
}

func TestGetStatsPropagatesRepositoryError(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)
	t.Parallel()

	repo := &fakeRepository{statsErr: errors.New("clickhouse unavailable")}

	got, err := New(repo).GetStats(context.Background())

	require.Error(err)
	assert.ErrorContains(err, "clickhouse unavailable")
	assert.Equal(model.Stats{}, got)
}

func TestGetStatsWithoutEventsReturnsEmptyStats(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)
	t.Parallel()

	got, err := New(&fakeRepository{}).GetStats(context.Background())

	require.NoError(err)
	assert.Zero(got.TotalRequests)
	assert.Zero(got.MostUsed.Hits)
	assert.Zero(got.MostUsed.Share)
}
