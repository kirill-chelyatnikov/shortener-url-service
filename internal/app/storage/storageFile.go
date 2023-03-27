package storage

import (
	"encoding/json"
	"fmt"
	"github.com/kirill-chelyatnikov/shortener-url-service/internal/app/models"
	"github.com/kirill-chelyatnikov/shortener-url-service/internal/config"
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"sync"
)

type storageFile struct {
	log     *logrus.Logger
	mutex   sync.RWMutex
	file    *os.File
	encoder *json.Encoder
	decoder *json.Decoder
}

// AddURL - функция записи данных в storage (file)
func (s *storageFile) AddURL(link *models.Link) error {
	if err := s.encoder.Encode(link); err != nil {
		return fmt.Errorf("can't encode data to add it to file, err: %s", err)
	}

	return nil
}

// GetURLByID - функция получения записи из storage (file)
func (s *storageFile) GetURLByID(id string) (string, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	f, err := os.Open(os.Getenv("FILE_STORAGE_PATH"))
	if err != nil {
		log.Fatalf("cant't create file storage, err: %s", err)
	}

	s.decoder = json.NewDecoder(f)
	var link models.Link

	for s.decoder.More() {
		s.decoder.Decode(&link)
		if link.Id == id {
			return link.BaseURL, nil
		}
	}

	return "", fmt.Errorf("can't find URL by id: %s", id)
}

func (s *storageFile) Close() error {
	err := s.file.Close()
	if err != nil {
		return err
	}

	return nil
}

func NewStorageFile(log *logrus.Logger, cfg *config.Config) *storageFile {
	f, err := os.OpenFile(os.Getenv("FILE_STORAGE_PATH"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		log.Fatalf("cant't create file storage, err: %s", err)
	}

	return &storageFile{
		log:     log,
		file:    f,
		encoder: json.NewEncoder(f),
	}
}

func (s *storageFile) CheckAndCloseFile(cfg *config.Config) {
	if cfg.App.FileStorage != "" {
		s.file.Close()
	}
}
