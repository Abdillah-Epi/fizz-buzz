package server

import (
	"net/http"

	"github.com/Abdillah-Epi/fizz-buzz/internal/config"
	"github.com/Abdillah-Epi/fizz-buzz/internal/fizzbuzz/handler"
	"github.com/Abdillah-Epi/fizz-buzz/internal/fizzbuzz/service"
	"github.com/Abdillah-Epi/fizz-buzz/internal/router"
)

type Server struct {
	Router http.Handler
}

func New(cfg config.Config) *Server {
	fizzBuzzService := service.NewFizzBuzzService(cfg)
	fizzBuzzHandler := handler.NewFizzBuzzHandler(
		fizzBuzzService,
	)

	appRouter := router.New(
		fizzBuzzHandler,
	)

	return &Server{
		Router: appRouter,
	}
}
