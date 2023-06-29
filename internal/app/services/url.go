package services

import (
	"context"
	"errors"
	"github.com/kirill-chelyatnikov/shortener-url-service/internal/app/models"
)

// Add - функция сервиса для добавления записи
func (s *ServiceURL) Add(ctx context.Context, link *models.Link) error {
	if len(link.ID) == 0 || len(link.BaseURL) == 0 || len(link.Hash) == 0 {
		return errors.New("empty url received")
	}

	if err := s.repository.AddURL(ctx, link); err != nil {
		return err
	}

	return nil
}

func (s *ServiceURL) AddBatch(ctx context.Context, links []*models.Link) error {
	if len(links) == 0 {
		return errors.New("passed an empty array of references")
	}

	if err := s.repository.AddURLSBatch(ctx, links); err != nil {
		return err
	}

	return nil
}

// Get - функция сервиса для получение записи по ID
func (s *ServiceURL) Get(ctx context.Context, id string) (string, error) {
	return s.repository.GetURLByID(ctx, id)
}

func (s *ServiceURL) GetAll(ctx context.Context, hash string) ([]*models.Link, error) {
	return s.repository.GetAllURLSByHash(ctx, hash)
}
