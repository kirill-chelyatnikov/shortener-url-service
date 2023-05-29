package storage

import (
	"errors"
	"fmt"
	"github.com/kirill-chelyatnikov/shortener-url-service/internal/app/models"
	"github.com/sirupsen/logrus"
	"sync"
)

type MapStorage struct {
	log   *logrus.Logger
	mutex sync.RWMutex
	data  map[string]models.Link
}

// AddURL - функция записи данных в storage (map)
func (s *MapStorage) AddURL(link *models.Link) error {
	s.mutex.Lock()

	s.data[link.ID] = models.Link{
		ID:      link.ID,
		BaseURL: link.BaseURL,
		Hash:    link.Hash,
	}
	s.mutex.Unlock()
	s.log.Infof("success write to map storage: id - %s, value - %s", link.ID, link.BaseURL)

	return nil
}

// GetURLByID - функция получения записи из storage (map)
func (s *MapStorage) GetURLByID(id string) (string, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if _, ok := s.data[id]; !ok {
		return "", fmt.Errorf("can't find URL by id: %s", id)
	}

	return s.data[id].BaseURL, nil
}

func (s *MapStorage) GetAllURLSByHash(hash string) ([]*models.Link, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	var links []*models.Link
	for _, v := range s.data {
		if v.Hash == hash {
			links = append(links, &v)
		}
	}

	if len(links) == 0 {
		return nil, errors.New("the user has no previously created links")
	}

	return links, nil
}

// Close - функция-заглушка для удовлетворения интерфейсу RepositoryInterface
func (s *MapStorage) Close() error {
	return nil
}

func NewMapStorage(log *logrus.Logger) *MapStorage {
	return &MapStorage{
		log:  log,
		data: make(map[string]models.Link),
	}
}
