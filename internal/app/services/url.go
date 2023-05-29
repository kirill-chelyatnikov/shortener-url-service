package services

import (
	"errors"
	"github.com/kirill-chelyatnikov/shortener-url-service/internal/app/models"
)

// Add - функция сервиса для добавления записи
func (s *ServiceURL) Add(link *models.Link) error {
	if len(link.ID) == 0 || len(link.BaseURL) == 0 || len(link.Hash) == 0 {
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

func (s *ServiceURL) GetAll(hash string) ([]*models.Link, error) {
	return s.repository.GetAllURLSByHash(hash)
}
