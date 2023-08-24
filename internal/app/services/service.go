package services

import (
	"context"
	"go.uber.org/zap"

	"github.com/kirill-chelyatnikov/shortener-url-service/internal/app/models"
	"github.com/kirill-chelyatnikov/shortener-url-service/internal/config"
)

type ServiceURL struct {
	log        *zap.SugaredLogger
	cfg        *config.Config
	repository RepositoryInterface
}

type RepositoryInterface interface {
	AddURL(ctx context.Context, link *models.Link) error
	AddURLSBatch(ctx context.Context, links []*models.Link) error
	GetURLByID(ctx context.Context, id string) (string, error)
	GetAllURLSByHash(ctx context.Context, hash string) ([]*models.Link, error)
	CheckBaseURLExist(ctx context.Context, link *models.Link) (bool, error)
	UpdateHash(ctx context.Context, link *models.Link) error
	Close() error
}

func NewServiceURL(log *zap.SugaredLogger, cfg *config.Config, repository RepositoryInterface) *ServiceURL {
	return &ServiceURL{
		log:        log,
		cfg:        cfg,
		repository: repository,
	}
}
