package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	handler "github.com/Abdillah-Epi/fizz-buzz/internal/fizzbuzz/handler"
)

func registerRoutes(r chi.Router, fizzBuzzHandler *handler.FizzBuzzHandler) {
	r.Get("/health", healthHandler)

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/fizzbuzz", func(r chi.Router) {
			fizzbuzzRoutes(r, fizzBuzzHandler)
		})
	})
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
