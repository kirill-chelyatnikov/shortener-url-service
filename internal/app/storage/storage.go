package storage

import (
	"fmt"
)

type storageMap struct {
	data map[string]string
}

func (s storageMap) AddURL(shortURL, baseURL string) {
	s.data[shortURL] = baseURL
}

func (s storageMap) GetURLByID(id string) (string, error) {
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
