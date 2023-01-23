package utils

import (
	"math/rand"
)

const (
	chars       = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ="
	charsSize   = len(chars)
	MaxSizeOfID = 6
)

func GenerateID() string {
	ID := ""

	for i := 0; i < MaxSizeOfID; i += 1 {
		index := rand.Int63n(int64(charsSize - 1))
		letter := chars[index]

		ID += string(letter)
	}

	return ID
}
