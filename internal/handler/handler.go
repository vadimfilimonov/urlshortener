package handler

import (
	"encoding/json"
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

type Input struct {
	Url string `json:"url"`
}

type Output struct {
	Result string `json:"result"`
}

func NewShorten(data storage.Data) func(http.ResponseWriter, *http.Request) {
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
				shortURL := fmt.Sprintf("%s/%s", Host, id)

				input := Input{}
				err = json.Unmarshal([]byte(body), &input)

				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				output, err := json.Marshal(Output{
					Result: shortURL,
				})

				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				data.Add(input.Url, id)
				w.Header().Add("Content-Type", "application/json")
				w.WriteHeader(http.StatusCreated)
				w.Write([]byte(output))
			}
		}
	}
}
