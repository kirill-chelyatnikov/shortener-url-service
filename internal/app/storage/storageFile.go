package storage

import (
	"encoding/json"
	"fmt"
	"github.com/kirill-chelyatnikov/shortener-url-service/internal/app/models"
	"github.com/kirill-chelyatnikov/shortener-url-service/internal/config"
	"github.com/sirupsen/logrus"
	"os"
	"sync"
)

type FileStorage struct {
	log     *logrus.Logger
	mutex   sync.RWMutex
	file    *os.File
	encoder *json.Encoder
	decoder *json.Decoder
}

// AddURL - функция записи данных в storage (file)
func (s *FileStorage) AddURL(link *models.Link) error {
	s.mutex.Lock()
	if err := s.encoder.Encode(link); err != nil {
		return fmt.Errorf("can't encode data to add it to file, err: %s", err)
	}
	s.mutex.Unlock()
	s.log.Infof("success write to file storage: id - %s, value - %s", link.ID, link.BaseURL)

	return nil
}

// GetURLByID - функция получения записи из storage (file)
func (s *FileStorage) GetURLByID(id string) (string, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	f, err := os.Open(s.file.Name())
	if err != nil {
		s.log.Fatalf("cant't open file storage, err: %s", err)
	}

	s.decoder = json.NewDecoder(f)
	var link models.Link

	for s.decoder.More() {
		err = s.decoder.Decode(&link)
		if err != nil {
			s.log.Fatalf("can't decode link, err: %s", err)
		}
		if link.ID == id {
			return link.BaseURL, nil
		}
	}

	return "", fmt.Errorf("can't find URL by id: %s", id)
}
func (s *FileStorage) GetAllURLSByHash(hash string) ([]*models.Link, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	//links := make([]models.Link, 0)
	//
	//for s

	return nil, nil
}

func (s *FileStorage) Close() error {
	err := s.file.Close()
	if err != nil {
		return err
	}

	return nil
}

func NewFileStorage(log *logrus.Logger, cfg *config.Config) *FileStorage {
	f, err := os.OpenFile(cfg.App.FileStorage, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		log.Fatalf("cant't create file storage, err: %s", err)
	}

	return &FileStorage{
		log:     log,
		file:    f,
		encoder: json.NewEncoder(f),
	}
}
