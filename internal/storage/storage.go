package storage

import (
	"github.com/VadimFilimonov/urlshortener/internal/config"
)

type Data interface {
	Get(shortenURL string) (string, error)
	GetItemsOfUser(userID string) ([]item, error)
	Add(originalURL, userID string) (shortenURL string, err error)
}

type item struct {
	userID      string
	ShortenURL  string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

func GetStorage(config config.Config) (Data, error) {
	if config.DatabaseDNS != "" {
		err := RunMigrations(config.DatabaseDNS)

		if err != nil {
			return nil, err
		}

		return NewDB(config.DatabaseDNS), nil
	}

	if config.FileStoragePath != "" {
		return NewFile(config.FileStoragePath), nil
	}

	return NewMemory(), nil
}
