package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Abdillah-Epi/fizz-buzz/internal/fizzbuzz/dto"
)

func getQueryParams(r *http.Request) (dto.FizzBuzzRequest, error) {
	query := r.URL.Query()

	int1, err := strconv.Atoi(query.Get("int1"))
	if err != nil {
		return dto.FizzBuzzRequest{}, fmt.Errorf("invalid int1")
	}

	int2, err := strconv.Atoi(query.Get("int2"))
	if err != nil {
		return dto.FizzBuzzRequest{}, fmt.Errorf("invalid int2")
	}

	limit, err := strconv.Atoi(query.Get("limit"))
	if err != nil {
		return dto.FizzBuzzRequest{}, fmt.Errorf("invalid limit")
	}

	str1 := query.Get("str1")
	str2 := query.Get("str2")

	request := dto.FizzBuzzRequest{
		Int1:  int1,
		Int2:  int2,
		Limit: limit,
		Str1:  str1,
		Str2:  str2,
	}
	return request, nil
}
