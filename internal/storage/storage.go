package storage

import "errors"

type DataItem struct {
	originalURL string
	shortURL    string
}

type Data struct {
	list map[string]DataItem
}

func New() Data {
	return Data{
		list: map[string]DataItem{},
	}
}

func (d Data) Get(id string) (string, error) {
	dataItem, ok := d.list[id]

	if !ok {
		return "", errors.New("incorrect id")
	}

	return dataItem.originalURL, nil
}

func (d Data) Add(originalURL, shortURL string) bool {
	d.list[shortURL] = DataItem{
		shortURL:    shortURL,
		originalURL: originalURL,
	}

	return true
}
