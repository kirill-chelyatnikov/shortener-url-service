package services

import (
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
	AddURL(link *models.Link) error
	GetURLByID(id string) (string, error)
	GetAllURLSByHash(hash string) ([]*models.Link, error)
	Close() error
}

func NewServiceURL(log *logrus.Logger, cfg *config.Config, repository RepositoryInterface) *ServiceURL {
	return &ServiceURL{
		log:        log,
		cfg:        cfg,
		repository: repository,
	}
}
