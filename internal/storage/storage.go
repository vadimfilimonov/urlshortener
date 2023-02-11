package storage

import "errors"

type Data struct {
	URLs map[string]string
}

func New() Data {
	return Data{
		URLs: map[string]string{},
	}
}

func (d Data) Get(shortenURL string) (string, error) {
	originalURL, ok := d.URLs[shortenURL]

	if !ok {
		return "", errors.New("incorrect shortenURL")
	}

	return originalURL, nil
}

func (d Data) Add(originalURL, shortenURL string) bool {
	d.URLs[shortenURL] = originalURL

	return true
}
