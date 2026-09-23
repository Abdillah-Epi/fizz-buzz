package router

import (
	"github.com/Abdillah-Epi/fizz-buzz/internal/analytics/handler"
	"github.com/go-chi/chi/v5"
)

func statsRoutes(r chi.Router, h *handler.StatsHandler) {
	r.Get("/", h.GetStats)
}
