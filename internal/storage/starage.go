package storage

import (
	"errors"
	"fmt"
)

type storageMap struct {
	data map[string]string
}

func (s storageMap) AddURL(shortURL, baseURL string) error {
	s.data[shortURL] = baseURL

	if _, ok := s.data[shortURL]; !ok {
		return errors.New("failed to write url to storage")
	}

	return nil
}

func (s storageMap) GetURLByID(id string) (string, error) {
	if _, ok := s.data[id]; !ok {
		return "", fmt.Errorf("can't find url by id: %s", id)
	}

	return s.data[id], nil
}

func NewStorage() *storageMap {
	return &storageMap{
		data: make(map[string]string),
	}
}
