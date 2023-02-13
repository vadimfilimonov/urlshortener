package main

import (
	"compress/gzip"
	"flag"
	"io"
	"net/http"
	"strings"

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
	r.Use(compressMiddleware)
	data := storage.New(config.FileStoragePath)

	r.Get("/{shortURL}", handler.New(data, config.BaseURL))
	r.Post("/", handler.New(data, config.BaseURL))
	r.Post("/api/shorten", handler.NewShorten(data, config.BaseURL))
	err = http.ListenAndServe(config.ServerAddress, r)

	if err != nil {
		panic(err)
	}
}

type writer struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w writer) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func compressMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
		if err != nil {
			io.WriteString(w, err.Error())
			return
		}
		defer gz.Close()

		w.Header().Set("Content-Encoding", "gzip")
		next.ServeHTTP(writer{ResponseWriter: w, Writer: gz}, r)
	})
}
