package storage

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

type Data struct {
	filename string
	URLs     map[string]string
}

func New(filename string) Data {
	if filename != "" {
		file, _ := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		file.Close()
	}

	return Data{
		filename: filename,
		URLs:     map[string]string{},
	}
}

func (d Data) Get(shortenURL string) (string, error) {
	if d.filename == "" {
		originalURL, ok := d.URLs[shortenURL]

		if !ok {
			return "", errors.New("incorrect shortenURL")
		}

		return originalURL, nil
	}

	data, err := os.ReadFile(d.filename)

	if err != nil {
		return "", err
	}

	rows := strings.Split(string(data), "\n")
	var originalURL string

	for _, row := range rows {
		if strings.Contains(row, shortenURL) {
			urls := strings.Split(row, " ")
			originalURL = urls[1]
			break
		}
	}

	if originalURL == "" {
		return "", errors.New("incorrect shortenURL")
	}

	return originalURL, nil
}

func (d Data) Add(originalURL, shortenURL string) bool {
	shouldSaveURLsToMemory := d.filename == ""

	if shouldSaveURLsToMemory {
		d.URLs[shortenURL] = originalURL
		return true
	}

	file, err := os.OpenFile(d.filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		return false
	}
	writer := bufio.NewWriter(file)
	data := fmt.Sprintf("%s %s\n", shortenURL, originalURL)
	_, err = writer.Write([]byte(data))

	if err != nil {
		return false
	}

	err = writer.Flush()
	file.Close()

	return err == nil
}
