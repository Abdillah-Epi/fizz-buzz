package model

import "time"

type RequestEvent struct {
	Timestamp  time.Time
	Int1       uint32
	Int2       uint32
	LimitValue uint32
	Str1       string
	Str2       string
}
