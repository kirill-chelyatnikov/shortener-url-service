package storage

import (
	"context"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"sync"

	"github.com/kirill-chelyatnikov/shortener-url-service/internal/app/models"
)

type MapStorage struct {
	log   *zap.SugaredLogger
	mutex sync.RWMutex
	data  map[string]*models.Link
}

// AddURL - функция записи данных в storage (map)
func (s *MapStorage) AddURL(ctx context.Context, link *models.Link) error {
	s.mutex.Lock()

	s.data[link.ID] = &models.Link{
		ID:      link.ID,
		BaseURL: link.BaseURL,
		Hash:    link.Hash,
	}
	s.mutex.Unlock()
	s.log.Infof("success write to map storage: id - %s, value - %s", link.ID, link.BaseURL)

	return nil
}

// GetURLByID - функция получения записи из storage (map)
func (s *MapStorage) GetURLByID(ctx context.Context, id string) (string, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if _, ok := s.data[id]; !ok {
		return "", fmt.Errorf("can't find URL by id: %s", id)
	}

	return s.data[id].BaseURL, nil
}

// GetAllURLSByHash - функция получения всех записей по хешу из storage (map)
func (s *MapStorage) GetAllURLSByHash(ctx context.Context, hash string) ([]*models.Link, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	var links []*models.Link
	for _, v := range s.data {
		if v.Hash == hash {
			links = append(links, v)
		}
	}

	if len(links) == 0 {
		return nil, errors.New("the user has no previously created links")
	}

	return links, nil
}

// функции-заглушки для удовлетворения интерфейсу репозитория

func (s *MapStorage) AddURLSBatch(ctx context.Context, links []*models.Link) error {
	return nil
}

func (s *MapStorage) CheckBaseURLExist(ctx context.Context, link *models.Link) (bool, error) {
	return false, nil
}

func (s *MapStorage) UpdateHash(ctx context.Context, link *models.Link) error {
	return nil
}

func (s *MapStorage) Close() error {
	return nil
}

func NewMapStorage(log *zap.SugaredLogger) *MapStorage {
	return &MapStorage{
		log:  log,
		data: make(map[string]*models.Link),
	}
}
