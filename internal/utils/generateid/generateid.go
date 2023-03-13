package utils

import (
	"math/rand"
	"time"
)

const (
	chars       = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ="
	charsSize   = len(chars)
	MaxSizeOfID = 6
)

func GenerateID() string {
	ID := ""
	randomGenerator := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < MaxSizeOfID; i += 1 {
		index := randomGenerator.Int63n(int64(charsSize - 1))
		letter := chars[index]

		ID += string(letter)
	}

	return ID
}
