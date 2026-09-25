package service

import (
	"fmt"
	"math"
	"strconv"

	config "github.com/Abdillah-Epi/fizz-buzz/internal/config"
	model "github.com/Abdillah-Epi/fizz-buzz/internal/fizzbuzz/model"
)

const MaxDivisor = int64(math.MaxUint32)

type FizzBuzzService struct {
	MaxLimit int
}

func NewFizzBuzzService(cfg config.Config) *FizzBuzzService {
	return &FizzBuzzService{
		MaxLimit: cfg.FizzBuzz.MaxLimit,
	}
}

func (s *FizzBuzzService) Generate(req model.Request) (model.Result, error) {
	if req.Int1 <= 0 {
		return model.Result{}, fmt.Errorf("int1 must be greater than 0")
	}

	if int64(req.Int1) > MaxDivisor {
		return model.Result{}, fmt.Errorf("int1 must be less than or equal to %d", MaxDivisor)
	}

	if req.Int2 <= 0 {
		return model.Result{}, fmt.Errorf("int2 must be greater than 0")
	}

	if int64(req.Int2) > MaxDivisor {
		return model.Result{}, fmt.Errorf("int2 must be less than or equal to %d", MaxDivisor)
	}

	if req.Limit <= 0 {
		return model.Result{}, fmt.Errorf("limit must be greater than 0")
	}

	if req.Limit > s.MaxLimit {
		return model.Result{}, fmt.Errorf("limit must be less than or equal to %d", s.MaxLimit)
	}

	if req.Str1 == "" {
		return model.Result{}, fmt.Errorf("str1 must not be empty")
	}

	if req.Str2 == "" {
		return model.Result{}, fmt.Errorf("str2 must not be empty")
	}

	values := fizzBuzz(req.Int1, req.Int2, req.Limit, req.Str1, req.Str2)

	return model.Result{
		Values: values,
	}, nil
}

func fizzBuzz(int1, int2, limit int, str1, str2 string) []string {
	values := make([]string, 0, limit)

	for i := 1; i <= limit; i++ {
		switch {
		case i%int1 == 0 && i%int2 == 0:
			values = append(values, str1+str2)

		case i%int1 == 0:
			values = append(values, str1)

		case i%int2 == 0:
			values = append(values, str2)

		default:
			values = append(values, strconv.Itoa(i))
		}
	}
	return values
}
