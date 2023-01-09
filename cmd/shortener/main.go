package main

import (
	"net/http"

	"github.com/VadimFilimonov/urlshortener/internal/handler"
	"github.com/VadimFilimonov/urlshortener/internal/storage"
)

func main() {
	http.HandleFunc("/", handler.New(storage.New()))
	http.ListenAndServe(":8080", nil)
}
