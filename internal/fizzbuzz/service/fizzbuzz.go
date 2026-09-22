package service

import (
	"strconv"

	model "github.com/Abdillah-Epi/fizz-buzz/internal/fizzbuzz/model"
)

type FizzBuzzService struct{}

func NewFizzBuzzService() *FizzBuzzService {
	return &FizzBuzzService{}
}

func (s *FizzBuzzService) Generate(req model.Request) model.Result {

	values := fizzBuzz(req.Int1, req.Int2, req.Limit, req.Str1, req.Str2)

	return model.Result{
		Values: values,
	}
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
