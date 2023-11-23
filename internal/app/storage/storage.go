package storage

import (
	"context"
	"go.uber.org/zap"

	"github.com/kirill-chelyatnikov/shortener-url-service/internal/app/services"
	"github.com/kirill-chelyatnikov/shortener-url-service/internal/config"
)

// NewStorage - функция получения хранилища в зафисимости от выбранного способа хранить ссылки (Map / File / DB)
func NewStorage(ctx context.Context, log *zap.SugaredLogger, cfg *config.Config) services.RepositoryInterface {
	if cfg.DB.CDN != "" {
		return NewPostgreSQLStorage(ctx, log, cfg)
	}

	if cfg.App.FileStorage != "" {
		return NewFileStorage(log, cfg)
	}

	return NewMapStorage(log)
}
