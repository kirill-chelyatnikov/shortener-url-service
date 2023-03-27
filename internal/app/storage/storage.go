package storage

import (
	"github.com/kirill-chelyatnikov/shortener-url-service/internal/app/services"
	"github.com/kirill-chelyatnikov/shortener-url-service/internal/config"
	"github.com/sirupsen/logrus"
)

func NewStorage(log *logrus.Logger, cfg *config.Config) services.RepositoryInterface {
	if cfg.App.FileStorage != "" {
		return NewStorageFile(log, cfg)
	}

	return NewStorageMap(log)
}
