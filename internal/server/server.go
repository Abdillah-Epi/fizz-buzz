package server

import (
	"fmt"
	"log/slog"
	"net/http"

	analyticsHandler "github.com/Abdillah-Epi/fizz-buzz/internal/analytics/handler"
	"github.com/Abdillah-Epi/fizz-buzz/internal/analytics/repository"
	analyticSservice "github.com/Abdillah-Epi/fizz-buzz/internal/analytics/service"
	"github.com/Abdillah-Epi/fizz-buzz/internal/config"
	"github.com/Abdillah-Epi/fizz-buzz/internal/fizzbuzz/handler"
	"github.com/Abdillah-Epi/fizz-buzz/internal/fizzbuzz/service"
	"github.com/Abdillah-Epi/fizz-buzz/internal/infrastructure/clickhouse"
	"github.com/Abdillah-Epi/fizz-buzz/internal/router"
)

type Server struct {
	Router     http.Handler
	ClickHouse *clickhouse.Client
}

func New(cfg config.Config, logger *slog.Logger) (*Server, error) {
	clickhouseClient, err := clickhouse.New(cfg.ClickHouse)
	if err != nil {
		return nil, fmt.Errorf("create clickhouse client: %w", err)
	}

	analyticsRepository := repository.NewClickHouseRepository(clickhouseClient)

	analyticsService := analyticSservice.New(analyticsRepository)

	statsHandler := analyticsHandler.NewStatsHandler(analyticsService)

	fizzBuzzService := service.NewFizzBuzzService(cfg)
	fizzBuzzHandler := handler.NewFizzBuzzHandler(fizzBuzzService, analyticsService, logger)

	handlers := router.RouterHandlers{
		FizzBuzzHandler: fizzBuzzHandler,
		StatsHandler:    statsHandler,
	}
	appRouter := router.New(handlers)

	return &Server{
		Router:     appRouter,
		ClickHouse: clickhouseClient,
	}, nil
}

func (s *Server) Close() error {
	if err := s.ClickHouse.Close(); err != nil {
		return fmt.Errorf("close clickhouse: %w", err)
	}

	return nil
}
