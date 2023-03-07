package storage

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
