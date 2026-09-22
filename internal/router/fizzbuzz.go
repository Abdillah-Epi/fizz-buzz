package router

import (
	handler "github.com/Abdillah-Epi/fizz-buzz/internal/fizzbuzz/handler"
	"github.com/go-chi/chi/v5"
)

func fizzbuzzRoutes(r chi.Router, fizzBuzzHandler *handler.FizzBuzzHandler) {
	r.Get("/", fizzBuzzHandler.Generate)
}
