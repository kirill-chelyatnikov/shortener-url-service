package storage

import (
	"github.com/kirill-chelyatnikov/shortener-url-service/internal/app/services"
	"github.com/kirill-chelyatnikov/shortener-url-service/internal/config"
	"github.com/sirupsen/logrus"
)

// NewStorage - функция получения хранилища в зафисимости от выбранного способа хранить ссылки (Map / File)
func NewStorage(log *logrus.Logger, cfg *config.Config) services.RepositoryInterface {
	if cfg.App.FileStorage != "" {
		return NewFileStorage(log)
	}

	return NewMapStorage(log)
}
