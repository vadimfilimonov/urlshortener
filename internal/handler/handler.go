package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/VadimFilimonov/urlshortener/internal/storage"
	utils "github.com/VadimFilimonov/urlshortener/internal/utils/generateid"
	"github.com/go-chi/chi/v5"
)

type Config struct {
	ServerAddress string `env:"SERVER_ADDRESS" envDefault:"localhost:8080"`
	BaseURL       string `env:"BASE_URL"`
}

func New(data storage.Data, config Config) func(http.ResponseWriter, *http.Request) {
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
				pathParts := make([]string, 0)
				pathParts = append(pathParts, config.ServerAddress)
				if config.BaseURL != "" {
					pathParts = append(pathParts, config.BaseURL)
				}
				pathParts = append(pathParts, id)
				shortURL := "http://" + strings.Join(pathParts, "/")

				data.Add(string(body), id)

				w.WriteHeader(http.StatusCreated)
				w.Write([]byte(shortURL))
			}
		}
	}
}

type Input struct {
	URL string `json:"url"`
}

type Output struct {
	Result string `json:"result"`
}

func NewShorten(data storage.Data, config Config) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			{
				body, err := io.ReadAll(r.Body)

				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				id := utils.GenerateID()
				pathParts := make([]string, 0)
				pathParts = append(pathParts, config.ServerAddress)
				if config.BaseURL != "" {
					pathParts = append(pathParts, config.BaseURL)
				}
				pathParts = append(pathParts, id)
				shortURL := "http://" + strings.Join(pathParts, "/")

				input := Input{}
				err = json.Unmarshal([]byte(body), &input)

				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}

				output, err := json.Marshal(Output{
					Result: shortURL,
				})

				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				data.Add(input.URL, id)
				w.Header().Add("Content-Type", "application/json")
				w.WriteHeader(http.StatusCreated)
				w.Write([]byte(output))
			}
		}
	}
}
