package main

import (
	"net/http"

	env "github.com/caarlos0/env/v6"

	"github.com/VadimFilimonov/urlshortener/internal/handler"
	"github.com/VadimFilimonov/urlshortener/internal/storage"
	"github.com/go-chi/chi/v5"
)

func main() {
	var config handler.Config
	err := env.Parse(&config)

	if err != nil {
		panic(err)
	}

	r := chi.NewRouter()
	data := storage.New()

	r.Get("/{shortURL}", handler.New(data, config))
	r.Post("/", handler.New(data, config))
	r.Post("/api/shorten", handler.NewShorten(data, config))
	err = http.ListenAndServe(config.ServerAddress, r)

	if err != nil {
		panic(err)
	}
}
