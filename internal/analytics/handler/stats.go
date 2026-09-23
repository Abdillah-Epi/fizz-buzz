package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Abdillah-Epi/fizz-buzz/internal/analytics/dto"
	"github.com/Abdillah-Epi/fizz-buzz/internal/analytics/service"
)

type StatsHandler struct {
	service *service.Service
}

func NewStatsHandler(service *service.Service) *StatsHandler {
	return &StatsHandler{
		service: service,
	}
}

func (h *StatsHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	stats, err := h.service.GetStats(r.Context())
	if err != nil {
		http.Error(w, "failed to get statistics", http.StatusInternalServerError)
		return
	}

	response := dto.StatsResponse{
		TotalRequests:        stats.TotalRequests,
		RequestsLastHour:     stats.RequestsLastHour,
		RequestsLast24Hours:  stats.RequestsLast24Hours,
		UniqueConfigurations: stats.UniqueConfigurations,
		MostUsed: dto.MostUsedRequest{
			Int1:        stats.MostUsed.Int1,
			Int2:        stats.MostUsed.Int2,
			Limit:       stats.MostUsed.LimitValue,
			Str1:        stats.MostUsed.Str1,
			Str2:        stats.MostUsed.Str2,
			Hits:        stats.MostUsed.Hits,
			Share:       stats.MostUsed.Share,
			FirstSeenAt: stats.MostUsed.FirstSeenAt,
			LastSeenAt:  stats.MostUsed.LastSeenAt,
		},
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}
