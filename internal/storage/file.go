package storage

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"golang.org/x/exp/slices"

	utils "github.com/VadimFilimonov/urlshortener/internal/utils/generateid"
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
	var status string

	for _, row := range rows {
		if strings.Contains(row, shortenURL) {
			columns := strings.Split(row, " ")
			originalURL = columns[1]
			status = columns[3]
			break
		}
	}

	if originalURL == "" {
		return "", errors.New("incorrect shortenURL")
	}

	if status == itemStatusDeleted {
		return "", URLHasBeenDeletedErr
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
		if row == "" {
			continue
		}

		columns := strings.Split(row, " ")
		shortenURL := columns[0]
		originalURL := columns[1]
		userIDOfOwner := columns[2]

		if userID != userIDOfOwner {
			continue
		}

		if slices.Contains(ids, shortenURL) {
			processedRow := fmt.Sprintf("%s %s %s %s", shortenURL, originalURL, userIDOfOwner, itemStatusDeleted)
			processedRows[index] = processedRow
		} else {
			processedRows[index] = row
		}
	}
	processedData := []byte(strings.Join(processedRows, "\n"))

	file, err := os.OpenFile(d.filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)

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
