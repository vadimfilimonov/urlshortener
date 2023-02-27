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

func New(data storage.Data, host string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		userIDCookie, ErrNoCookie := r.Cookie("userID")

		if ErrNoCookie != nil {
			userIDCookie = &http.Cookie{
				Name:  "userID",
				Value: utils.GenerateID(),
			}
		}

		http.SetCookie(w, userIDCookie)

		switch r.Method {
		case http.MethodGet:
			{
				shortenURL := chi.URLParam(r, "shortenURL")

				isURLEmpty := shortenURL == ""

				if isURLEmpty {
					http.Error(w, "shortenURL param is missed", http.StatusBadRequest)
					return
				}
				originalURL, err := data.Get(shortenURL)

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
				defer r.Body.Close()

				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				path := utils.GenerateID()
				shortenURL := fmt.Sprintf("%s/%s", host, path)

				data.Add(string(body), path, userIDCookie.Value)

				w.WriteHeader(http.StatusCreated)
				w.Write([]byte(shortenURL))
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

func NewShorten(data storage.Data, host string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		userIDCookie, ErrNoCookie := r.Cookie("userID")

		if ErrNoCookie != nil {
			userIDCookie = &http.Cookie{
				Name:  "userID",
				Value: utils.GenerateID(),
			}
		}

		http.SetCookie(w, userIDCookie)

		body, err := io.ReadAll(r.Body)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		path := utils.GenerateID()
		shortenURL := fmt.Sprintf("%s/%s", host, path)

		input := Input{}
		err = json.Unmarshal([]byte(body), &input)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		output, err := json.Marshal(Output{
			Result: shortenURL,
		})

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		data.Add(input.URL, path, userIDCookie.Value)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(output))
	}
}

func NewUserUrls(data storage.Data) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		userIDCookie, ErrNoCookie := r.Cookie("userID")

		if ErrNoCookie != nil {
			userIDCookie = &http.Cookie{
				Name:  "userID",
				Value: utils.GenerateID(),
			}
		}

		http.SetCookie(w, userIDCookie)

		items, err := data.GetItemsOfUser(userIDCookie.Value)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if len(items) == 0 {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		response, err := json.Marshal(items)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Add("Content-Type", "application/json")
		w.Write(response)
	}
}
