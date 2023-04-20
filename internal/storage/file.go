package storage

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	utils "github.com/VadimFilimonov/urlshortener/internal/utils/generateid"
	"golang.org/x/exp/slices"
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
				status:      columns[3],
			}
			items = append(items, item)
		}
	}

	return items, nil
}

func (d dataFile) Add(originalURL, userID string) (string, error) {
	file, err := os.OpenFile(d.filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		return "", err
	}
	writer := bufio.NewWriter(file)
	shortenURLPath := utils.GenerateID()
	data := fmt.Sprintf("%s %s %s %s\n", shortenURLPath, originalURL, userID, itemStatusCreated)
	_, err = writer.Write([]byte(data))

	if err != nil {
		return "", err
	}

	err = writer.Flush()
	file.Close()

	if err != nil {
		return "", err
	}

	return shortenURLPath, nil
}

func (d dataFile) Delete(ids []string, userID string) error {
	data, err := os.ReadFile(d.filename)

	if err != nil {
		return err
	}

	rows := strings.Split(string(data), "\n")
	processedRows := make([]string, len(rows))

	for index, row := range rows {
		columns := strings.Split(row, " ")
		shortenURL := columns[0]
		originalURL := columns[1]
		userID := columns[2]

		if slices.Contains(ids, shortenURL) {
			data := fmt.Sprintf("%s %s %s %s\n", shortenURL, originalURL, userID, itemStatusDeleted)
			processedRows[index] = data
		} else {
			processedRows[index] = row
		}
	}
	processedData := []byte(strings.Join(processedRows, "\n"))

	file, err := os.OpenFile(d.filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)

	if err != nil {
		return err
	}

	writer := bufio.NewWriter(file)
	_, err = writer.Write(processedData)

	if err != nil {
		return err
	}

	err = writer.Flush()

	if err != nil {
		return err
	}

	err = file.Close()

	if err != nil {
		return err
	}

	return nil
}
