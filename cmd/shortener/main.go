package main

import (
	"net/http"

	"github.com/VadimFilimonov/urlshortener/internal/handler"
)

func main() {
	http.HandleFunc("/", handler.New())
	http.ListenAndServe(":8080", nil)
}
