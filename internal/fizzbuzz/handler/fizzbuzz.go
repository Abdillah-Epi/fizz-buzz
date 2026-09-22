package handler

import (
	"encoding/json"
	"net/http"

	dto "github.com/Abdillah-Epi/fizz-buzz/internal/fizzbuzz/dto"
	model "github.com/Abdillah-Epi/fizz-buzz/internal/fizzbuzz/model"
	service "github.com/Abdillah-Epi/fizz-buzz/internal/fizzbuzz/service"
)

type FizzBuzzHandler struct {
	service *service.FizzBuzzService
}

func NewFizzBuzzHandler(service *service.FizzBuzzService) *FizzBuzzHandler {
	return &FizzBuzzHandler{
		service: service,
	}
}

func (h *FizzBuzzHandler) Generate(w http.ResponseWriter, r *http.Request) {

	request, err := getQueryParams(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	result := h.service.Generate(model.Request{
		Int1:  request.Int1,
		Int2:  request.Int2,
		Limit: request.Limit,
		Str1:  request.Str1,
		Str2:  request.Str2,
	})

	response := dto.FizzBuzzResponse{
		Values: result.Values,
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}
