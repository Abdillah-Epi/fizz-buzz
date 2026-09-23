package router

import (
	"net/http"
	"time"

	analyticsHandler "github.com/Abdillah-Epi/fizz-buzz/internal/analytics/handler"
	fizzbuzzHandler "github.com/Abdillah-Epi/fizz-buzz/internal/fizzbuzz/handler"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type RouterHandlers struct {
	FizzBuzzHandler *fizzbuzzHandler.FizzBuzzHandler
	StatsHandler    *analyticsHandler.StatsHandler
}

func New(handlers RouterHandlers) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.ClientIPFromRemoteAddr)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	registerRoutes(r, handlers)

	return r
}
