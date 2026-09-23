package service

import (
	"context"

	"github.com/Abdillah-Epi/fizz-buzz/internal/analytics/model"
	"github.com/Abdillah-Epi/fizz-buzz/internal/analytics/repository"
)

type Service struct {
	repository repository.Repository
}

func New(repository repository.Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) RecordEvent(ctx context.Context, event model.RequestEvent) error {
	return s.repository.RecordEvent(ctx, event)
}

func (s *Service) GetStats(ctx context.Context) (model.Stats, error) {
	return s.repository.GetStats(ctx)
}
