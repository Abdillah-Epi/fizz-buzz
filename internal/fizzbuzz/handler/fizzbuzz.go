package handler

import (
	"net/http"
)

func FizzBuzzHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
