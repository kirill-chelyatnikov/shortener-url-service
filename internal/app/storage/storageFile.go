package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"os"
	"sync"

	"github.com/kirill-chelyatnikov/shortener-url-service/internal/app/models"
	"github.com/kirill-chelyatnikov/shortener-url-service/internal/config"
)

type FileStorage struct {
	log     *zap.SugaredLogger
	mutex   sync.RWMutex
	file    *os.File
	encoder *json.Encoder
	decoder *json.Decoder
}

// AddURL - функция записи данных в storage (file)
func (s *FileStorage) AddURL(ctx context.Context, link *models.Link) error {
	s.mutex.Lock()
	if err := s.encoder.Encode(link); err != nil {
		return fmt.Errorf("can't encode data to add it to file, err: %s", err)
	}
	s.mutex.Unlock()
	s.log.Infof("success write to file storage: id - %s, value - %s", link.ID, link.BaseURL)

	return nil
}

// GetURLByID - функция получения записи из storage (file)
func (s *FileStorage) GetURLByID(ctx context.Context, id string) (*models.Link, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	f, err := os.Open(s.file.Name())
	if err != nil {
		s.log.Fatalf("cant't open file storage, err: %s", err)
	}

	s.decoder = json.NewDecoder(f)
	var link *models.Link

	for s.decoder.More() {
		err = s.decoder.Decode(&link)
		if err != nil {
			s.log.Fatalf("can't decode link, err: %s", err)
		}
		if link.ID == id {
			return link, nil
		}
	}

	return nil, fmt.Errorf("can't find URL by id: %s", id)
}

// функции-заглушки для удовлетворения интерфейсу репозитория

func (s *FileStorage) GetAllURLSByHash(_ context.Context, hash string) ([]*models.Link, error) {
	return nil, nil
}

func (s *FileStorage) AddURLSBatch(_ context.Context, links []*models.Link) error {
	return nil
}

func (s *FileStorage) DeleteURLSBatch(_ context.Context, links []string) error {
	return nil
}

func (s *FileStorage) CheckBaseURLExist(_ context.Context, link *models.Link) (bool, error) {
	return false, nil
}

func (s *FileStorage) UpdateHash(_ context.Context, link *models.Link) error {
	return nil
}

func (s *FileStorage) Close() error {
	err := s.file.Close()
	if err != nil {
		return err
	}

	return nil
}

func NewFileStorage(log *zap.SugaredLogger, cfg *config.Config) *FileStorage {
	f, err := os.OpenFile(cfg.App.FileStorage, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0o777)
	if err != nil {
		log.Fatalf("cant't create file storage, err: %s", err)
	}

	return &FileStorage{
		log:     log,
		file:    f,
		encoder: json.NewEncoder(f),
	}
}
