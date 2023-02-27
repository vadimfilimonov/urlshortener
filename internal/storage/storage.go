package storage

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

type item struct {
	userID      string
	ShortenURL  string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type Data struct {
	filename string
	items    map[string]item
}

func New(filename string) Data {
	if filename != "" {
		file, _ := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		file.Close()
	}

	return Data{
		filename: filename,
		items:    map[string]item{},
	}
}

func (d Data) Get(shortenURL string) (string, error) {
	if d.filename == "" {
		item, ok := d.items[shortenURL]

		if !ok {
			return "", errors.New("incorrect shortenURL")
		}

		return item.OriginalURL, nil
	}

	data, err := os.ReadFile(d.filename)

	if err != nil {
		return "", err
	}

	rows := strings.Split(string(data), "\n")
	var originalURL string

	for _, row := range rows {
		if strings.Contains(row, shortenURL) {
			columns := strings.Split(row, " ")
			originalURL = columns[1]
			break
		}
	}

	if originalURL == "" {
		return "", errors.New("incorrect shortenURL")
	}

	return originalURL, nil
}

func (d Data) GetItemsOfUser(userID string) ([]item, error) {
	items := make([]item, 0)

	if d.filename == "" {
		for _, item := range d.items {
			if item.userID == userID {
				items = append(items, item)
			}
		}

		return items, nil
	}

	data, err := os.ReadFile(d.filename)

	if err != nil {
		return nil, err
	}

	rows := strings.Split(string(data), "\n")

	for _, row := range rows {
		if strings.Contains(row, userID) {
			columns := strings.Split(row, " ")
			item := item{
				userID:      columns[2],
				ShortenURL:  columns[0],
				OriginalURL: columns[1],
			}
			items = append(items, item)
		}
	}

	return items, nil
}

func (d Data) Add(originalURL, shortenURL, userID string) bool {
	shouldSaveURLsToMemory := d.filename == ""

	if shouldSaveURLsToMemory {
		d.items[shortenURL] = item{
			userID:      userID,
			ShortenURL:  shortenURL,
			OriginalURL: originalURL,
		}
		return true
	}

	file, err := os.OpenFile(d.filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		return false
	}
	writer := bufio.NewWriter(file)
	data := fmt.Sprintf("%s %s %s\n", shortenURL, originalURL, userID)
	_, err = writer.Write([]byte(data))

	if err != nil {
		return false
	}

	err = writer.Flush()
	file.Close()

	return err == nil
}
