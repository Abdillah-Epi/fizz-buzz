package dto

import "time"

type StatsResponse struct {
	TotalRequests        uint64          `json:"total_requests"`
	RequestsLastHour     uint64          `json:"requests_last_hour"`
	RequestsLast24Hours  uint64          `json:"requests_last_24h"`
	UniqueConfigurations uint64          `json:"unique_configurations"`
	MostUsed             MostUsedRequest `json:"most_used"`
}

type MostUsedRequest struct {
	Int1        uint32    `json:"int1"`
	Int2        uint32    `json:"int2"`
	Limit       uint32    `json:"limit"`
	Str1        string    `json:"str1"`
	Str2        string    `json:"str2"`
	Hits        uint64    `json:"hits"`
	Share       float64   `json:"share"`
	FirstSeenAt time.Time `json:"first_seen_at"`
	LastSeenAt  time.Time `json:"last_seen_at"`
}
