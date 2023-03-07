package storage

import (
	"errors"

	utils "github.com/VadimFilimonov/urlshortener/internal/utils/generateid"
)

type memoryItems map[string]item

func NewMemory() memoryItems {
	return memoryItems{}
}

func (items memoryItems) Get(shortenURL string) (string, error) {
	item, ok := items[shortenURL]

	if !ok {
		return "", errors.New("incorrect shortenURL")
	}

	return item.OriginalURL, nil
}

func (items memoryItems) GetItemsOfUser(userID string) ([]item, error) {
	userItems := make([]item, 0)

	for _, item := range items {
		if item.userID == userID {
			userItems = append(userItems, item)
		}
	}

	return userItems, nil
}

func (items memoryItems) Add(originalURL, userID string) (string, error) {
	shortenURLPath := utils.GenerateID()

	items[shortenURLPath] = item{
		userID:      userID,
		ShortenURL:  shortenURLPath,
		OriginalURL: originalURL,
	}

	return shortenURLPath, nil
}
