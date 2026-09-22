package service

import (
	"fmt"
	"testing"

	"github.com/Abdillah-Epi/fizz-buzz/internal/config"
	"github.com/Abdillah-Epi/fizz-buzz/internal/fizzbuzz/model"
	"github.com/stretchr/testify/assert"
)

func NewConfig() config.Config {
	return config.Config{
		FizzBuzz: config.FizzBuzzConfig{
			MaxLimit: 10000,
		},
	}
}

func initFizzBuzzService(int1, int2, limit int, str1, str2 string) (*FizzBuzzService, model.Request) {
	cfg := NewConfig()
	fb := NewFizzBuzzService(cfg)
	req := model.Request{
		Int1:  int1,
		Int2:  int2,
		Str1:  str1,
		Str2:  str2,
		Limit: limit,
	}

	return fb, req
}

func TestFizzBuzz(t *testing.T) {
	assert := assert.New(t)
	expected := []string{"1", "2", "Fizz", "4", "Buzz", "Fizz", "7", "8", "Fizz", "Buzz", "11", "Fizz", "13", "14", "FizzBuzz"}

	res := fizzBuzz(3, 5, 15, "Fizz", "Buzz")

	for i, v := range res {
		if i > len(expected) {
			assert.True(false)
			break
		}
		assert.Equal(v, expected[i])
	}
}

func TestInt1MustBePositive(t *testing.T) {
	assert := assert.New(t)
	fb, req := initFizzBuzzService(0, 5, 15, "Fizz", "Buzz")
	res, err := fb.Generate(req)
	if err == nil {
		assert.True(false)
	}
	assert.Empty(res)
	assert.Equal("int1 must be greater than 0", err.Error())
}

func TestInt2MustBePositive(t *testing.T) {
	assert := assert.New(t)
	fb, req := initFizzBuzzService(3, 0, 15, "Fizz", "Buzz")
	res, err := fb.Generate(req)
	if err == nil {
		assert.True(false)
	}
	assert.Empty(res)
	assert.Equal("int2 must be greater than 0", err.Error())
}

func TestLimitMustBePositive(t *testing.T) {
	assert := assert.New(t)
	fb, req := initFizzBuzzService(3, 5, 0, "Fizz", "Buzz")
	res, err := fb.Generate(req)
	if err == nil {
		assert.True(false)
	}
	assert.Empty(res)
	assert.Equal("limit must be greater than 0", err.Error())
}

func TestLimitExceedingMaxLimit(t *testing.T) {
	assert := assert.New(t)
	fb, req := initFizzBuzzService(3, 5, 100000, "Fizz", "Buzz")
	res, err := fb.Generate(req)
	if err == nil {
		assert.True(false)
	}
	assert.Empty(res)
	expected := fmt.Sprintf("limit must be less than or equal to %v", fb.MaxLimit)
	assert.Equal(expected, err.Error())
}

func TestStr1MustNotBeEmpty(t *testing.T) {
	assert := assert.New(t)
	fb, req := initFizzBuzzService(3, 5, 15, "", "Buzz")
	res, err := fb.Generate(req)
	if err == nil {
		assert.True(false)
	}
	assert.Empty(res)
	assert.Equal("str1 must not be empty", err.Error())
}

func TestStr2MustNotBeEmpty(t *testing.T) {
	assert := assert.New(t)
	fb, req := initFizzBuzzService(3, 5, 15, "Fizz", "")
	res, err := fb.Generate(req)
	if err == nil {
		assert.True(false)
	}
	assert.Empty(res)
	assert.Equal("str2 must not be empty", err.Error())
}
