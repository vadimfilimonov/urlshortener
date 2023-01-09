package storage

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

func (d Data) Get(id string) (string, bool) {
	dataItem, ok := d.list[id]
	return dataItem.originalURL, ok
}

func (d Data) Add(originalURL, shortURL string) bool {
	d.list[shortURL] = DataItem{
		shortURL:    shortURL,
		originalURL: originalURL,
	}

	return true
}
