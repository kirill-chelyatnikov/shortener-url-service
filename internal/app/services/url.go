package services

import (
	"errors"
	"math/rand"
	"time"
)

// Add - функция сервиса для добавления записи
func (s *ServiceURL) Add(shortURL, baseURL string) error {
	if len(shortURL) == 0 || len(baseURL) == 0 {
		return errors.New("empty url received")
	}

	s.repository.AddURL(shortURL, baseURL)

	return nil
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
