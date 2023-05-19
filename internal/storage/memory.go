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

	if item.status == itemStatusDeleted {
		return "", URLHasBeenDeletedErr
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
		status:      itemStatusCreated,
	}

	return shortenURLPath, nil
}

func (items memoryItems) Delete(ids []string, userID string) error {
	for _, id := range ids {
		itemCopy := items[id]

		if itemCopy.userID == userID {
			itemCopy.status = itemStatusDeleted
			items[id] = itemCopy
		}
	}
	return nil
}
