package main

import (
	"bytes"
	"compress/gzip"
	"flag"
	"io"
	"log"
	"net/http"
	"strings"

	env "github.com/caarlos0/env/v6"

	"github.com/VadimFilimonov/urlshortener/internal/handler"
	"github.com/VadimFilimonov/urlshortener/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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
		log.Fatal(err)
	}

	serverAddress := flag.String("a", "localhost:8080", "адрес запуска HTTP-сервера")
	BaseURL := flag.String("b", "http://localhost:8080", "базовый адрес результирующего сокращённого URL")
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
	r.Use(decompressMiddleware)
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

func decompressMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		gz, err := gzip.NewReader(r.Body)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer gz.Close()

		body, err := io.ReadAll(gz)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		r.Body = io.NopCloser(bytes.NewBuffer(body))
		next.ServeHTTP(w, r)
	})
}
