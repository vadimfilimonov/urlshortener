package storage

type Data interface {
	Get(shortenURL string) (string, error)
	GetItemsOfUser(userId string) ([]item, error)
	Add(originalURL, shortenURL, userID string) bool
}

type item struct {
	userID      string
	ShortenURL  string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}
