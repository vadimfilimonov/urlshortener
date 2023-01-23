package handler

import (
	"fmt"
	"io"
	"net/http"

	"github.com/VadimFilimonov/urlshortener/internal/storage"
	utils "github.com/VadimFilimonov/urlshortener/internal/utils/generateid"
	"github.com/go-chi/chi/v5"
)

const (
	Host string = "http://localhost:8080"
)

func New(data storage.Data) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			{
				shortURL := chi.URLParam(r, "shortURL")

				isURLEmpty := shortURL == ""

				if isURLEmpty {
					http.Error(w, "shortURL param is missed", http.StatusBadRequest)
					return
				}
				originalURL, err := data.Get(shortURL)

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

				id := utils.GenerateID()
				shortURL := fmt.Sprintf("%s/%s", Host, id)

				data.Add(string(body), id)

				w.WriteHeader(http.StatusCreated)
				w.Write([]byte(shortURL))
			}
		}
	}
}
