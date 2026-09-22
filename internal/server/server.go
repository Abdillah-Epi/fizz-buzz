package server

import (
	"fmt"
	"net/http"

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

func New(cfg config.Config) (*Server, error) {
	clickhouseClient, err := clickhouse.New(cfg.ClickHouse)
	if err != nil {
		return nil, fmt.Errorf("create clickhouse client: %w", err)
	}

	fizzBuzzService := service.NewFizzBuzzService(cfg)
	fizzBuzzHandler := handler.NewFizzBuzzHandler(
		fizzBuzzService,
	)

	appRouter := router.New(
		fizzBuzzHandler,
	)

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
