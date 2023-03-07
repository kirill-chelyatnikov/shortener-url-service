package storage

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"sync"
)

type storageMap struct {
	log   *logrus.Logger
	mutex sync.RWMutex
	data  map[string]string
}

// AddURL - функция записи данных в storage
func (s *storageMap) AddURL(shortURL, baseURL string) {
	s.mutex.Lock()
	s.data[shortURL] = baseURL
	s.mutex.Unlock()
	s.log.Infof("success write to storage: key - %s, value - %s", shortURL, baseURL)
}

// GetURLByID - функция получения записи из storage
func (s *storageMap) GetURLByID(id string) (string, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if _, ok := s.data[id]; !ok {
		return "", fmt.Errorf("can't find URL by id: %s", id)
	}

	return s.data[id], nil
}

func NewStorage(log *logrus.Logger) *storageMap {
	return &storageMap{
		log:  log,
		data: make(map[string]string),
	}
}
