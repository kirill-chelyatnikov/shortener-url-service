package services

import (
	"math/rand"
	"time"
)

// Add - функция сервиса для добавления записи
func (s *ServiceURL) Add(shortURL, baseURL string) {
	s.repository.AddURL(shortURL, baseURL)
}

// Get - функция сервиса для получение записи по ID
func (s *ServiceURL) Get(id string) (string, error) {
	return s.repository.GetURLByID(id)
}

// GenerateShortURL - функция генерации короткого URL
func (s *ServiceURL) GenerateShortURL() string {
	rand.Seed(time.Now().UnixNano())
	var chars = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0987654321")
	str := make([]rune, s.cfg.App.ShortedURLLen)

	for i := range str {
		str[i] = chars[rand.Intn(len(chars))]
	}

	return string(str)
}
