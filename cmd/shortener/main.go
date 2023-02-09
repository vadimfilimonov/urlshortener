package main

import (
	"net/http"

	"github.com/VadimFilimonov/urlshortener/internal/handler"
	"github.com/VadimFilimonov/urlshortener/internal/storage"
	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()
	data := storage.New()
	r.Get("/{shortURL}", handler.New(data))
	r.Post("/", handler.New(data))
	r.Post("/api/shorten", handler.NewShorten(data))
	http.ListenAndServe(":8080", r)
}
