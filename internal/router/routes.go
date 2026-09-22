package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	handler "github.com/Abdillah-Epi/fizz-buzz/internal/fizzbuzz/handler"
	service "github.com/Abdillah-Epi/fizz-buzz/internal/fizzbuzz/service"
)

func registerRoutes(r chi.Router) {
	r.Get("/health", healthHandler)

	fizzBuzzService := service.NewFizzBuzzService()
	fizzBuzzHandler := handler.NewFizzBuzzHandler(fizzBuzzService)

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/fizzbuzz", func(r chi.Router) {
			fizzbuzzRoutes(r, fizzBuzzHandler)
		})
	})
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
