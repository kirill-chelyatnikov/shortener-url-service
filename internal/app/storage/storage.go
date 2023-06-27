package storage

import (
	"context"
	"github.com/kirill-chelyatnikov/shortener-url-service/internal/app/services"
	"github.com/kirill-chelyatnikov/shortener-url-service/internal/config"
	"github.com/sirupsen/logrus"
)

// NewStorage - функция получения хранилища в зафисимости от выбранного способа хранить ссылки (Map / File)
func NewStorage(ctx context.Context, log *logrus.Logger, cfg *config.Config) services.RepositoryInterface {
	if cfg.Db.CDN != "" {
		return NewPostgreSQLStorage(ctx, log, cfg)
	}

	if cfg.App.FileStorage != "" {
		return NewFileStorage(log, cfg)
	}

	return NewMapStorage(log)
}
