package encode

import (
	"math/rand"
)

const letters = "abcdefghijklmnopqrstuvwxyz"
const lettersSize = 26

const maxSizeOfEncodedText = 6

func Encode(text string) string {
	encodedText := ""

	for i := 0; i < maxSizeOfEncodedText; i += 1 {
		index := rand.Int63n(lettersSize - 1)
		letter := letters[index]

		encodedText += string(letter)
	}

	return encodedText
}
