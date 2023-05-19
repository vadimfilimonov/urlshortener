package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/VadimFilimonov/urlshortener/internal/constants"
	"github.com/VadimFilimonov/urlshortener/internal/storage"
	utils "github.com/VadimFilimonov/urlshortener/internal/utils/generateid"
)

func NewGet(data storage.Data, host string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		shortenURL := chi.URLParam(r, "shortenURL")

		if len(shortenURL) == 0 {
			http.Error(w, "shortenURL param is missed", http.StatusBadRequest)
			return
		}
		originalURL, err := data.Get(shortenURL)

		if errors.Is(err, storage.ErrURLHasBeenDeleted) {
			http.Error(w, err.Error(), http.StatusGone)
			return
		}

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Location", originalURL)
		w.WriteHeader(http.StatusTemporaryRedirect)
	}
}

func NewPost(data storage.Data, host string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		userIDCookieValue := manageUserIDCookie(w, r)

		body, err := io.ReadAll(r.Body)
		defer r.Body.Close()

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		shortenURLPath, err := data.Add(string(body), userIDCookieValue)
		shortenURL := fmt.Sprintf("%s/%s", host, shortenURLPath)

		if errors.Is(err, constants.ErrURLAlreadyExists) {
			w.WriteHeader(http.StatusConflict)
			w.Write([]byte(shortenURL))
			return
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(shortenURL))
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
		defer r.Body.Close()

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var requestBody ShortenInput
		err = json.Unmarshal([]byte(body), &requestBody)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		shortenURLPath, errDataAdd := data.Add(requestBody.URL, userIDCookieValue)
		shortenURL := fmt.Sprintf("%s/%s", host, shortenURLPath)

		responseJSON, err := json.Marshal(ShortenOutput{
			Result: shortenURL,
		})

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if errors.Is(errDataAdd, constants.ErrURLAlreadyExists) {
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusConflict)
			w.Write([]byte(responseJSON))
			return
		}
		if errDataAdd != nil {
			http.Error(w, errDataAdd.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(responseJSON))
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
		defer r.Body.Close()

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

func NewGetUserUrls(data storage.Data, host string) func(http.ResponseWriter, *http.Request) {
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

func NewDeleteUserUrls(data storage.Data) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		userIDCookieValue := manageUserIDCookie(w, r)
		body, err := io.ReadAll(r.Body)
		defer r.Body.Close()

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		ids := make([]string, 0)
		err = json.Unmarshal([]byte(body), &ids)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = data.Delete(ids, userIDCookieValue)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusAccepted)
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
		defer ctx.Done()
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
