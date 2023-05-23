package main

import (
	"bytes"
	"compress/gzip"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/VadimFilimonov/urlshortener/internal/config"
	"github.com/VadimFilimonov/urlshortener/internal/handler"
	"github.com/VadimFilimonov/urlshortener/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	config := config.New()
	config.Parse()

	data, err := storage.GetStorage(config)
	if err != nil {
		log.Fatal(err)
	}

	r := chi.NewRouter()
	r.Use(decompressMiddleware)
	r.Use(middleware.Compress(5))
	r.Get("/{shortenURL}", handler.NewGet(data, config.BaseURL))
	r.Post("/", handler.NewPost(data, config.BaseURL))
	r.Post("/api/shorten", handler.NewShorten(data, config.BaseURL))
	r.Post("/api/shorten/batch", handler.NewShortenBatch(data, config.BaseURL))
	r.Get("/api/user/urls", handler.NewGetUserUrls(data, config.BaseURL))
	r.Delete("/api/user/urls", handler.NewDeleteUserUrls(data))
	r.Get("/ping", handler.NewPing(config.DatabaseDNS))
	err = http.ListenAndServe(config.ServerAddress, r)

	if err != nil {
		log.Fatal(err)
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
