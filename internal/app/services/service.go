package services

import (
	"context"
	"github.com/kirill-chelyatnikov/shortener-url-service/internal/app/models"
	"github.com/kirill-chelyatnikov/shortener-url-service/internal/config"
	"github.com/sirupsen/logrus"
)

type ServiceURL struct {
	log        *logrus.Logger
	cfg        *config.Config
	repository RepositoryInterface
}

type RepositoryInterface interface {
	AddURL(ctx context.Context, link *models.Link) error
	GetURLByID(ctx context.Context, id string) (string, error)
	GetAllURLSByHash(ctx context.Context, hash string) ([]*models.Link, error)
	Close() error
}

func NewServiceURL(log *logrus.Logger, cfg *config.Config, repository RepositoryInterface) *ServiceURL {
	return &ServiceURL{
		log:        log,
		cfg:        cfg,
		repository: repository,
	}
}
