package main

import (
	"fmt"
	"net/http"

	env "github.com/caarlos0/env/v6"

	"github.com/VadimFilimonov/urlshortener/internal/handler"
	"github.com/VadimFilimonov/urlshortener/internal/storage"
	"github.com/go-chi/chi/v5"
)

type Config struct {
	ServerAddress string `env:"SERVER_ADDRESS"`
	BaseURL       string `env:"BASE_URL"`
}

func main() {
	config := Config{}
	err := env.Parse(&config)

	if err != nil {
		panic(err)
	}
	fmt.Println(config)

	r := chi.NewRouter()
	data := storage.New()

	r.Get("/{shortURL}", handler.New(data))
	r.Post("/", handler.New(data))
	r.Post("/api/shorten", handler.NewShorten(data))
	err = http.ListenAndServe(config.ServerAddress+":8080", r)

	if err != nil {
		panic(err)
	}
}
