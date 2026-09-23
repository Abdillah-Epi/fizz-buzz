package repository

import (
	"context"

	"github.com/Abdillah-Epi/fizz-buzz/internal/analytics/model"
)

type Repository interface {
	RecordEvent(ctx context.Context, event model.RequestEvent) error
	GetStats(ctx context.Context) (model.Stats, error)
}
