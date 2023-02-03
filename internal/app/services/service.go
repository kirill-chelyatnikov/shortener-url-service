package services

import (
	"github.com/kirill-chelyatnikov/shortener-url-service/internal/config"
	"github.com/sirupsen/logrus"
)

type ServiceURL struct {
	log        *logrus.Logger
	cfg        *config.Config
	repository RepositoryInterface
}

type RepositoryInterface interface {
	AddURL(shortURL, baseURL string)
	GetURLByID(id string) (string, error)
}

func NewServiceURL(log *logrus.Logger, cfg *config.Config, repository RepositoryInterface) *ServiceURL {
	return &ServiceURL{
		log:        log,
		cfg:        cfg,
		repository: repository,
	}
}
