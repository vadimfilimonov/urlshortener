package shortstring

import (
	"math/rand"
)

const (
	chars                = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ="
	charsSize            = int64(len(chars))
	maxSizeOfshortString = 6
)

func Generate() string {
	shortString := ""

	for i := 0; i < maxSizeOfshortString; i += 1 {
		index := rand.Int63n(charsSize - 1)
		letter := chars[index]

		shortString += string(letter)
	}

	return shortString
}
