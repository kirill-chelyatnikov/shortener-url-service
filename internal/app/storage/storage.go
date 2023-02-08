package storage

import (
	"fmt"
	"sync"
)

type storageMap struct {
	mutex sync.RWMutex
	data  map[string]string
}

// AddURL - функция записи данных в storage
func (s *storageMap) AddURL(shortURL, baseURL string) {
	s.mutex.Lock()
	s.data[shortURL] = baseURL
	s.mutex.Unlock()
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

func NewStorage() *storageMap {
	return &storageMap{
		data: make(map[string]string),
	}
}
