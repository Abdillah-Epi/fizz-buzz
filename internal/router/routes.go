package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func registerRoutes(r chi.Router, handlers RouterHandlers) {
	r.Get("/health", healthHandler)

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/fizzbuzz", func(r chi.Router) {
			fizzbuzzRoutes(r, handlers.FizzBuzzHandler)
		})

		r.Route("/stats", func(r chi.Router) {
			statsRoutes(r, handlers.StatsHandler)
		})
	})
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
