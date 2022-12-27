package main

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strings"
)

type DataItem struct {
	originalUrl string
	shortUrl    string
}

type Data struct {
	list map[string]DataItem
}

func (d Data) Get(id string) (DataItem, bool) {
	dataItem, ok := d.list[id]
	return dataItem, ok
}

func (d Data) Add(dataItem DataItem) bool {
	id := dataItem.shortUrl
	d.list[id] = dataItem

	return true
}

func Encode(text string) string {
	encodedText := fmt.Sprint(rand.Int63n(1000))

	return encodedText
}

var data = Data{
	list: map[string]DataItem{},
}

func Handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		{
			path := r.URL.Path
			isUrlEmpty := path == "/"
			if isUrlEmpty {
				http.Error(w, "Empty URL", http.StatusBadRequest)
				return
			}
			slices := strings.Split(path, "/")
			id := slices[len(slices)-1]
			dataItem, isDataItemExist := data.Get(id)

			if !isDataItemExist {
				http.Error(w, "Incorrect Id", http.StatusNotFound)
				return
			}

			w.Header().Set("Location", dataItem.originalUrl)
			w.WriteHeader(http.StatusTemporaryRedirect)
		}
	case http.MethodPost:
		{
			body, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			encodedUrl := Encode(string(body))
			data.Add(DataItem{
				originalUrl: string(body),
				shortUrl:    encodedUrl,
			})

			w.WriteHeader(http.StatusCreated)
			fmt.Println(encodedUrl)
			w.Write([]byte(encodedUrl))
		}
	}
}

func main() {
	http.HandleFunc("/", Handler)
	http.ListenAndServe(":8080", nil)
}

// POST / принимает в теле запроса строку URL для сокращения и возвращает ответ с сокращённым URL в виде текстовой строки в теле.
// GET /{id} принимает в качестве URL-параметра идентификатор сокращённого URL и возвращает ответ с оригинальным URL в HTTP-заголовке Location.
// Нужно учесть некорректные запросы и возвращать для них ответ с кодом 400.
