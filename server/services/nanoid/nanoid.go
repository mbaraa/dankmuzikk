package nanoid

import (
	"bytes"
	"math/rand"
	"time"
)

const defaultAlphabet = "_-0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

var random = rand.New(rand.NewSource(time.Now().UnixNano()))

// Generate generates a random nanoid with length of 12
func Generate() string {
	return GenerateWithLength(12)
}

// Generate generates a random nanoid with a custom length; 10 <= length <= 200
// If the size was out of that range, it defaults to 12.
func GenerateWithLength(size int) string {
	if size < 10 || size > 200 {
		size = 12
	}
	buf := bytes.NewBuffer([]byte{})
	for i := 0; i < size; i++ {
		buf.WriteByte(defaultAlphabet[random.Intn(len(defaultAlphabet))])
		random.Seed(time.Now().UnixNano())
	}
	return buf.String()
}
