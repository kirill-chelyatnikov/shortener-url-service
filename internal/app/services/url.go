package services

import (
	"errors"
	"github.com/kirill-chelyatnikov/shortener-url-service/internal/app/models"
	"math/rand"
	"time"
)

// Add - функция сервиса для добавления записи
func (s *ServiceURL) Add(link *models.Link) error {
	if len(link.ID) == 0 || len(link.BaseURL) == 0 {
		return errors.New("empty url received")
	}

	err := s.repository.AddURL(link)
	if err != nil {
		return err
	}

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
