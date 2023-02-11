package main

import (
	"net/http"

	env "github.com/caarlos0/env/v6"

	"github.com/VadimFilimonov/urlshortener/internal/handler"
	"github.com/VadimFilimonov/urlshortener/internal/storage"
	"github.com/go-chi/chi/v5"
)

const (
	Host string = "localhost:8080"
)

type Config struct {
	ServerAddress   string `env:"SERVER_ADDRESS"`
	BaseURL         string `env:"BASE_URL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
}

func main() {
	config := Config{
		ServerAddress: Host,
		BaseURL:       "http://" + Host,
	}
	err := env.Parse(&config)

	if err != nil {
		panic(err)
	}

	r := chi.NewRouter()
	data := storage.New()

	r.Get("/{shortURL}", handler.New(data, config.BaseURL))
	r.Post("/", handler.New(data, config.BaseURL))
	r.Post("/api/shorten", handler.NewShorten(data, config.BaseURL))
	err = http.ListenAndServe(config.ServerAddress, r)

	if err != nil {
		panic(err)
	}
}
