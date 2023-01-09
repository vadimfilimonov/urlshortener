package handler

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/VadimFilimonov/urlshortener/internal/utils/shortstring"
)

type DataItem struct {
	originalURL string
	shortURL    string
}

type Data struct {
	list map[string]DataItem
}

func (d Data) Get(id string) (DataItem, bool) {
	dataItem, ok := d.list[id]
	return dataItem, ok
}

func (d Data) Add(dataItem DataItem) bool {
	id := dataItem.shortURL
	d.list[id] = dataItem

	return true
}

var data = Data{
	list: map[string]DataItem{},
}

func Handler(w http.ResponseWriter, r *http.Request) {
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
			dataItem, isDataItemExist := data.Get(id)
			fmt.Println(path)

			if !isDataItemExist {
				http.Error(w, "Incorrect Id", http.StatusBadRequest)
				return
			}

			w.Header().Set("Location", dataItem.originalURL)
			w.WriteHeader(http.StatusTemporaryRedirect)
		}
	case http.MethodPost:
		{
			body, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			shortURL := shortstring.Generate()
			data.Add(DataItem{
				originalURL: string(body),
				shortURL:    shortURL,
			})
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte(shortURL))
		}
	}
}

// POST / принимает в теле запроса строку URL для сокращения и возвращает ответ с сокращённым URL в виде текстовой строки в теле.
// GET /{id} принимает в качестве URL-параметра идентификатор сокращённого URL и возвращает ответ с оригинальным URL в HTTP-заголовке Location.
// Нужно учесть некорректные запросы и возвращать для них ответ с кодом 400.
