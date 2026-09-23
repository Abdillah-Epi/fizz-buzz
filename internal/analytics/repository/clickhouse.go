package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Abdillah-Epi/fizz-buzz/internal/analytics/model"
	"github.com/Abdillah-Epi/fizz-buzz/internal/infrastructure/clickhouse"
)

type ClickHouseRepository struct {
	db *clickhouse.Client
}

func NewClickHouseRepository(db *clickhouse.Client) *ClickHouseRepository {
	return &ClickHouseRepository{
		db: db,
	}
}

func (r *ClickHouseRepository) RecordEvent(ctx context.Context, event model.RequestEvent) error {
	const query = `
		INSERT INTO request_events
		(
			timestamp,
			int1,
			int2,
			limit_value,
			str1,
			str2
		)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	err := r.db.Conn().Exec(ctx, query, event.Timestamp, event.Int1, event.Int2, event.LimitValue, event.Str1, event.Str2)
	if err != nil {
		return fmt.Errorf("insert request event: %w", err)
	}
	return nil
}

func (r *ClickHouseRepository) GetStats(ctx context.Context) (model.Stats, error) {
	const query = `
		WITH
			total_requests AS (
				SELECT count() AS total
				FROM request_events
			),
			configuration_stats AS (
				SELECT
					int1,
					int2,
					limit_value,
					str1,
					str2,
					count() AS hits,
					min(timestamp) AS first_seen_at,
					max(timestamp) AS last_seen_at
				FROM request_events
				GROUP BY
					int1,
					int2,
					limit_value,
					str1,
					str2
				ORDER BY hits DESC
				LIMIT 1
			)
		SELECT
			total_requests.total,

			(
				SELECT count()
				FROM request_events
				WHERE timestamp >= now() - INTERVAL 1 HOUR
			) AS requests_last_hour,

			(
				SELECT count()
				FROM request_events
				WHERE timestamp >= now() - INTERVAL 24 HOUR
			) AS requests_last_24h,

			(
				SELECT count()
				FROM (
					SELECT
						int1,
						int2,
						limit_value,
						str1,
						str2
					FROM request_events
					GROUP BY
						int1,
						int2,
						limit_value,
						str1,
						str2
				)
			) AS unique_configurations,

			configuration_stats.int1,
			configuration_stats.int2,
			configuration_stats.limit_value,
			configuration_stats.str1,
			configuration_stats.str2,
			configuration_stats.hits,

			configuration_stats.first_seen_at,
			configuration_stats.last_seen_at

		FROM total_requests
		CROSS JOIN configuration_stats
	`

	var stats model.Stats

	var (
		totalRequests        uint64
		requestsLastHour     uint64
		requestsLast24Hours  uint64
		uniqueConfigurations uint64

		int1       uint32
		int2       uint32
		limitValue uint32
		str1       string
		str2       string
		hits       uint64
		firstSeen  time.Time
		lastSeen   time.Time
	)

	err := r.db.Conn().QueryRow(ctx, query).Scan(
		&totalRequests,
		&requestsLastHour,
		&requestsLast24Hours,
		&uniqueConfigurations,
		&int1,
		&int2,
		&limitValue,
		&str1,
		&str2,
		&hits,
		&firstSeen,
		&lastSeen,
	)
	if errors.Is(err, sql.ErrNoRows) {
		// No events recorded yet: the aggregate query returns no rows. Report an
		// empty dataset instead of an error so /stats can answer with zeros.
		return model.Stats{}, nil
	}
	if err != nil {
		return model.Stats{}, fmt.Errorf("get analytics stats: %w", err)
	}

	var share float64
	if totalRequests > 0 {
		share = float64(hits) / float64(totalRequests) * 100
	}

	stats = model.Stats{
		TotalRequests:        totalRequests,
		RequestsLastHour:     requestsLastHour,
		RequestsLast24Hours:  requestsLast24Hours,
		UniqueConfigurations: uniqueConfigurations,
		MostUsed: model.MostUsedRequest{
			Int1:        int1,
			Int2:        int2,
			LimitValue:  limitValue,
			Str1:        str1,
			Str2:        str2,
			Hits:        hits,
			Share:       share,
			FirstSeenAt: firstSeen,
			LastSeenAt:  lastSeen,
		},
	}

	return stats, nil
}
