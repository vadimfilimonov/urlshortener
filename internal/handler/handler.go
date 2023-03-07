package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"database/sql"

	"github.com/VadimFilimonov/urlshortener/internal/constants"
	"github.com/VadimFilimonov/urlshortener/internal/storage"
	utils "github.com/VadimFilimonov/urlshortener/internal/utils/generateid"
	"github.com/go-chi/chi/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func New(data storage.Data, host string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		userIDCookieValue := manageUserIDCookie(w, r)

		switch r.Method {
		case http.MethodGet:
			{
				shortenURL := chi.URLParam(r, "shortenURL")

				if len(shortenURL) == 0 {
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

				shortenURLPath, errDataAdd := data.Add(string(body), userIDCookieValue)
				shortenURL := fmt.Sprintf("%s/%s", host, shortenURLPath)

				if errors.Is(errDataAdd, constants.ErrURLAlreadyExists) {
					w.WriteHeader(http.StatusConflict)
					w.Write([]byte(shortenURL))
					return
				} else if errDataAdd != nil {
					http.Error(w, errDataAdd.Error(), http.StatusInternalServerError)
					return
				}

				w.WriteHeader(http.StatusCreated)
				w.Write([]byte(shortenURL))
			}
		}
	}
}

type ShortenInput struct {
	URL string `json:"url"`
}

type ShortenOutput struct {
	Result string `json:"result"`
}

func NewShorten(data storage.Data, host string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		userIDCookieValue := manageUserIDCookie(w, r)

		body, err := io.ReadAll(r.Body)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		input := ShortenInput{}
		err = json.Unmarshal([]byte(body), &input)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		shortenURLPath, errDataAdd := data.Add(input.URL, userIDCookieValue)
		shortenURL := fmt.Sprintf("%s/%s", host, shortenURLPath)

		output, err := json.Marshal(ShortenOutput{
			Result: shortenURL,
		})

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if errors.Is(errDataAdd, constants.ErrURLAlreadyExists) {
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusConflict)
			w.Write([]byte(output))
			return
		} else if errDataAdd != nil {
			http.Error(w, errDataAdd.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(output))
	}
}

type ShortenBatchInputItem struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

type ShortenBatchOutputItem struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

func NewShortenBatch(data storage.Data, host string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		userIDCookieValue := manageUserIDCookie(w, r)

		body, err := io.ReadAll(r.Body)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		input := make([]ShortenBatchInputItem, 0)
		err = json.Unmarshal([]byte(body), &input)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		outputList := make([]ShortenBatchOutputItem, len(input))

		for i, item := range input {
			shortenURLPath, err := data.Add(item.OriginalURL, userIDCookieValue)
			shortenURL := fmt.Sprintf("%s/%s", host, shortenURLPath)

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			outputList[i] = ShortenBatchOutputItem{
				CorrelationID: item.CorrelationID,
				ShortURL:      shortenURL,
			}
		}

		output, err := json.Marshal(outputList)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(output))
	}
}

type URLData = struct {
	ShortenURL  string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

func NewUserUrls(data storage.Data, host string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		userIDCookieValue := manageUserIDCookie(w, r)

		items, err := data.GetItemsOfUser(userIDCookieValue)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if len(items) == 0 {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		normalizedItems := make([]URLData, len(items))
		for index, item := range items {
			normalizedItems[index] = URLData{
				ShortenURL:  fmt.Sprintf("%s/%s", host, item.ShortenURL),
				OriginalURL: item.OriginalURL,
			}
		}
		response, err := json.Marshal(normalizedItems)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Add("Content-Type", "application/json")
		w.Write(response)
	}
}

func NewPing(DBPath string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		manageUserIDCookie(w, r)

		if DBPath == "" {
			http.Error(w, "empty path to database", http.StatusInternalServerError)
			return
		}

		db, err := sql.Open("pgx", DBPath)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer db.Close()

		ctx, cancel := context.WithTimeout(r.Context(), time.Second*1)
		defer cancel()
		err = db.PingContext(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		w.WriteHeader(http.StatusOK)
	}
}

func manageUserIDCookie(w http.ResponseWriter, r *http.Request) string {
	userIDCookie, ErrNoCookie := r.Cookie("userID")

	if ErrNoCookie != nil {
		userIDCookie = &http.Cookie{
			Name:  "userID",
			Value: utils.GenerateID(),
		}
	}

	http.SetCookie(w, userIDCookie)

	return userIDCookie.Value
}
