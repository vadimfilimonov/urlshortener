package config

import (
	"flag"
	"log"

	env "github.com/caarlos0/env/v6"
)

type Config struct {
	ServerAddress   string `env:"SERVER_ADDRESS"`
	BaseURL         string `env:"BASE_URL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
}

func New() Config {
	return Config{}
}

func (c *Config) Parse() {
	err := env.Parse(c)

	if err != nil {
		log.Fatal(err)
	}

	serverAddress := flag.String("a", "localhost:8080", "адрес запуска HTTP-сервера")
	BaseURL := flag.String("b", "http://localhost:8080", "базовый адрес результирующего сокращённого URL")
	FileStoragePath := flag.String("f", "", "путь до файла с сокращёнными URL")
	flag.Parse()

	if c.ServerAddress == "" {
		c.ServerAddress = *serverAddress
	}

	if c.BaseURL == "" {
		c.BaseURL = *BaseURL
	}

	if c.FileStoragePath == "" {
		c.FileStoragePath = *FileStoragePath
	}
}
