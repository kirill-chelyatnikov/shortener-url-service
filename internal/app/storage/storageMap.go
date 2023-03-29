package storage

import (
	"fmt"
	"github.com/kirill-chelyatnikov/shortener-url-service/internal/app/models"
	"github.com/sirupsen/logrus"
	"sync"
)

type MapStorage struct {
	log   *logrus.Logger
	mutex sync.RWMutex
	data  map[string]string
}

// AddURL - функция записи данных в storage (map)
func (s *MapStorage) AddURL(link *models.Link) error {
	s.mutex.Lock()
	s.data[link.Id] = link.BaseURL
	s.mutex.Unlock()
	s.log.Infof("success write to map storage: id - %s, value - %s", link.Id, link.BaseURL)

	return nil
}

// GetURLByID - функция получения записи из storage (map)
func (s *MapStorage) GetURLByID(id string) (string, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if _, ok := s.data[id]; !ok {
		return "", fmt.Errorf("can't find URL by id: %s", id)
	}

	return s.data[id], nil
}

// Close - функция-заглушка для удовлетворения интерфейсу RepositoryInterface
func (s *MapStorage) Close() error {
	return nil
}

func NewMapStorage(log *logrus.Logger) *MapStorage {
	return &MapStorage{
		log:  log,
		data: make(map[string]string),
	}
}
