package pkg

import (
	"math/rand"
	"time"
)

// GenerateRandomString - функция генерации короткого URL/cookie
func GenerateRandomString() string {
	rand.Seed(time.Now().UnixNano())
	var chars = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0987654321")
	str := make([]rune, 10)

	for i := range str {
		str[i] = chars[rand.Intn(len(chars))]
	}

	return string(str)
}
