package storage

import (
	"errors"

	"github.com/VadimFilimonov/urlshortener/internal/config"
)

type Data interface {
	Get(shortenURL string) (string, error)
	GetItemsOfUser(userID string) ([]item, error)
	Add(originalURL, userID string) (shortenURL string, err error)
	Delete(ids []string, userID string) error
}

type item struct {
	userID      string
	ShortenURL  string `json:"short_url"`
	OriginalURL string `json:"original_url"`
	status      string
}

const (
	itemStatusCreated = "created"
	itemStatusDeleted = "deleted"
)

var ErrURLHasBeenDeleted = errors.New("url has been deleted")

func GetStorage(config config.Config) (Data, error) {
	if config.DatabaseDNS != "" {
		db, err := InitDB(config.DatabaseDNS)

		if err != nil {
			return nil, err
		}

		return NewDB(db), nil
	}

	if config.FileStoragePath != "" {
		return NewFile(config.FileStoragePath), nil
	}

	return NewMemory(), nil
}
