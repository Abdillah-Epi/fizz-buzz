package handler

import (
	"log/slog"
	"net/http"
	"time"

	analyticsModel "github.com/Abdillah-Epi/fizz-buzz/internal/analytics/model"
	analyticsService "github.com/Abdillah-Epi/fizz-buzz/internal/analytics/service"
	dto "github.com/Abdillah-Epi/fizz-buzz/internal/fizzbuzz/dto"
	model "github.com/Abdillah-Epi/fizz-buzz/internal/fizzbuzz/model"
	service "github.com/Abdillah-Epi/fizz-buzz/internal/fizzbuzz/service"
	apphttp "github.com/Abdillah-Epi/fizz-buzz/internal/http"
)

type FizzBuzzHandler struct {
	service          *service.FizzBuzzService
	analyticsService *analyticsService.Service
	logger           *slog.Logger
}

func NewFizzBuzzHandler(service *service.FizzBuzzService, analyticsService *analyticsService.Service, logger *slog.Logger) *FizzBuzzHandler {
	return &FizzBuzzHandler{
		service:          service,
		analyticsService: analyticsService,
		logger:           logger,
	}
}

func (h *FizzBuzzHandler) Generate(w http.ResponseWriter, r *http.Request) {
	log := h.logger
	req, err := getQueryParams(r)
	if err != nil {
		apphttp.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.service.Generate(model.Request{
		Int1:  req.Int1,
		Int2:  req.Int2,
		Limit: req.Limit,
		Str1:  req.Str1,
		Str2:  req.Str2,
	})
	if err != nil {
		apphttp.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	event := analyticsModel.RequestEvent{
		Timestamp:  time.Now().UTC(),
		Int1:       uint32(req.Int1),
		Int2:       uint32(req.Int2),
		LimitValue: uint32(req.Limit),
		Str1:       req.Str1,
		Str2:       req.Str2,
	}

	if err := h.analyticsService.RecordEvent(r.Context(), event); err != nil {
		log.Error("failed to record analytics event", "error", err)
	}

	response := dto.FizzBuzzResponse{
		Values: result.Values,
	}

	apphttp.SendResponse(w, http.StatusOK, response)
}
