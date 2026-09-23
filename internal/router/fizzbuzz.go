package router

import (
	"github.com/Abdillah-Epi/fizz-buzz/internal/fizzbuzz/handler"
	"github.com/go-chi/chi/v5"
)

func fizzbuzzRoutes(r chi.Router, h *handler.FizzBuzzHandler) {
	r.Get("/", h.Generate)
}
