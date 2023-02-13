package main

import (
	"flag"
	"net/http"

	env "github.com/caarlos0/env/v6"

	"github.com/VadimFilimonov/urlshortener/internal/handler"
	"github.com/VadimFilimonov/urlshortener/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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
	config := Config{}
	err := env.Parse(&config)

	if err != nil {
		panic(err)
	}

	serverAddress := flag.String("a", Host, "адрес запуска HTTP-сервера")
	BaseURL := flag.String("b", "http://"+Host, "базовый адрес результирующего сокращённого URL")
	FileStoragePath := flag.String("f", "", "путь до файла с сокращёнными URL")
	flag.Parse()

	if config.ServerAddress == "" {
		config.ServerAddress = *serverAddress
	}

	if config.BaseURL == "" {
		config.BaseURL = *BaseURL
	}

	if config.FileStoragePath == "" {
		config.FileStoragePath = *FileStoragePath
	}

	r := chi.NewRouter()
	r.Use(middleware.Compress(5))
	data := storage.New(config.FileStoragePath)

	r.Get("/{shortURL}", handler.New(data, config.BaseURL))
	r.Post("/", handler.New(data, config.BaseURL))
	r.Post("/api/shorten", handler.NewShorten(data, config.BaseURL))
	err = http.ListenAndServe(config.ServerAddress, r)

	if err != nil {
		panic(err)
	}
}
