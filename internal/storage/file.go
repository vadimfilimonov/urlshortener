package storage

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

type dataFile struct {
	filename string
	items    map[string]item
}

func NewFile(filename string) dataFile {
	return dataFile{
		filename: filename,
		items:    map[string]item{},
	}
}

func (d dataFile) Get(shortenURL string) (string, error) {
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

func (d dataFile) GetItemsOfUser(userID string) ([]item, error) {
	items := make([]item, 0)
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

func (d dataFile) Add(originalURL, shortenURL, userID string) error {
	file, err := os.OpenFile(d.filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		return err
	}
	writer := bufio.NewWriter(file)
	data := fmt.Sprintf("%s %s %s\n", shortenURL, originalURL, userID)
	_, err = writer.Write([]byte(data))

	if err != nil {
		return err
	}

	err = writer.Flush()
	file.Close()

	if err != nil {
		return err
	}

	return nil
}
