package router

import (
	"net/http"
	"time"

	"github.com/Abdillah-Epi/fizz-buzz/internal/fizzbuzz/handler"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func New(fizzBuzzHandler *handler.FizzBuzzHandler) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.ClientIPFromRemoteAddr)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	registerRoutes(r, fizzBuzzHandler)

	return r
}
