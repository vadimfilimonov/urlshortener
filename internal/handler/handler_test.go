package handler

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/VadimFilimonov/urlshortener/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	Host string = "http://localhost:8080"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name       string
		request    string
		body       string
		method     string
		statusCode int
	}{
		{
			name:       "Shorten url generated",
			request:    Host,
			body:       "https://filimonovvadim.t.me",
			method:     http.MethodPost,
			statusCode: http.StatusCreated,
		},
		{
			name:       "Empty relative url",
			request:    Host,
			body:       "",
			method:     http.MethodGet,
			statusCode: http.StatusBadRequest,
		},
		{
			name:       "Invalid relative url",
			request:    fmt.Sprintf("%s/hash", Host),
			body:       "",
			method:     http.MethodGet,
			statusCode: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := strings.NewReader(tt.body)
			request := httptest.NewRequest(tt.method, tt.request, body)
			w := httptest.NewRecorder()
			h := http.HandlerFunc(New(storage.New(""), tt.request))
			h.ServeHTTP(w, request)

			result := w.Result()
			assert.Equal(t, tt.statusCode, result.StatusCode)
			bodyResult, err := io.ReadAll(result.Body)
			require.NoError(t, err)
			err = result.Body.Close()
			require.NoError(t, err)
			assert.NotEmpty(t, string(bodyResult))
		})
	}
}

func TestNewShorten(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		statusCode int
	}{
		{
			name:       "Shorten url generated",
			body:       "{\"url\":\"https://filimonovvadim.t.me\"}",
			statusCode: http.StatusCreated,
		},
		{
			name:       "Empty relative url",
			body:       "",
			statusCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := strings.NewReader(tt.body)
			request := httptest.NewRequest(http.MethodPost, fmt.Sprintf("%s/api/shorten", Host), body)
			w := httptest.NewRecorder()
			h := http.HandlerFunc(NewShorten(storage.New(""), Host))
			h.ServeHTTP(w, request)

			result := w.Result()
			assert.Equal(t, tt.statusCode, result.StatusCode)
			bodyResult, err := io.ReadAll(result.Body)
			require.NoError(t, err)
			err = result.Body.Close()
			require.NoError(t, err)
			assert.NotEmpty(t, string(bodyResult))
		})
	}
}
