package services

import (
	"context"
	"errors"
	"github.com/kirill-chelyatnikov/shortener-url-service/internal/app/models"
	"github.com/kirill-chelyatnikov/shortener-url-service/pkg"
)

// Add - функция сервиса для добавления/изменения записи
func (s *ServiceURL) Add(ctx context.Context, link *models.Link) (bool, error) {
	// Проверка на пустоту переданных полей
	if len(link.BaseURL) == 0 || len(link.Hash) == 0 {
		return false, errors.New("empty url received")
	}

	// Проверяем наличие URL в БД
	check, err := s.repository.CheckBaseURLExist(ctx, link)
	if err != nil {
		return false, err
	}

	/* Если URL в базе уже существует, то добавляем хеш пользователя в массив хэшей данной записи,
	также передаем информацию в хендлер о том, что запись мы проапдейтили, а не вставили новую, нужно для понимания
	какой код ответа следует отдать */
	if check {
		if err = s.repository.UpdateHash(ctx, link); err != nil {
			return false, err
		}
		return true, nil
	} else {
		/* Если URL не найден, то просто пишем его в БД.
		Присваиваем записи ID (в случае апдейта записи, ID мы берем из БД) */
		link.ID = pkg.GenerateRandomString()
		if err = s.repository.AddURL(ctx, link); err != nil {
			return false, err
		}

		return false, nil
	}
}

// AddBatch - функция сервиса для добавления записей "пачкой"
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

// GetAll - функция сервиса для получения всех записей по хешу
func (s *ServiceURL) GetAll(ctx context.Context, hash string) ([]*models.Link, error) {
	return s.repository.GetAllURLSByHash(ctx, hash)
}
