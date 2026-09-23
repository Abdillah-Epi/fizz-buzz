package model

import "time"

type Stats struct {
	TotalRequests        uint64
	RequestsLastHour     uint64
	RequestsLast24Hours  uint64
	UniqueConfigurations uint64
	MostUsed             MostUsedRequest
}

type MostUsedRequest struct {
	Int1        uint32
	Int2        uint32
	LimitValue  uint32
	Str1        string
	Str2        string
	Hits        uint64
	Share       float64
	FirstSeenAt time.Time
	LastSeenAt  time.Time
}
