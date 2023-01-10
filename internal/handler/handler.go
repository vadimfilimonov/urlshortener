package handler

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/VadimFilimonov/urlshortener/internal/storage"
	"github.com/VadimFilimonov/urlshortener/internal/utils/shortstring"
)

func New(data storage.Data) func(http.ResponseWriter, *http.Request) {
	h := func(w http.ResponseWriter, r *http.Request) {
		Handler(w, r, data)
	}

	return h
}

func Handler(w http.ResponseWriter, r *http.Request, data storage.Data) {
	switch r.Method {
	case http.MethodGet:
		{
			path := r.URL.Path
			isURLEmpty := path == "/"

			if isURLEmpty {
				http.Error(w, "Empty URL", http.StatusBadRequest)
				return
			}

			slices := strings.Split(path, "/")
			id := slices[len(slices)-1]
			originalURL, err := data.Get(id)

			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			w.Header().Set("Location", originalURL)
			w.WriteHeader(http.StatusTemporaryRedirect)
		}
	case http.MethodPost:
		{
			body, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			id := shortstring.Generate()
			shortURL := fmt.Sprintf("http://localhost:8080/%s", id)

			data.Add(string(body), id)

			w.WriteHeader(http.StatusCreated)
			w.Write([]byte(shortURL))
		}
	}
}

// POST / принимает в теле запроса строку URL для сокращения и возвращает ответ с сокращённым URL в виде текстовой строки в теле.
// GET /{id} принимает в качестве URL-параметра идентификатор сокращённого URL и возвращает ответ с оригинальным URL в HTTP-заголовке Location.
// Нужно учесть некорректные запросы и возвращать для них ответ с кодом 400.
