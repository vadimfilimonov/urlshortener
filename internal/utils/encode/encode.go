package encode

import (
	"fmt"
	"math/rand"
)

func Encode(text string) string {
	encodedText := fmt.Sprint(rand.Int63n(1000))

	return encodedText
}
