package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
