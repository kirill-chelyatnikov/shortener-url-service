package server

import (
	"github.com/go-chi/chi"
	"github.com/kirill-chelyatnikov/shortener-url-service/internal/config"
	"github.com/sirupsen/logrus"
	"net/http"
)

// HTTPServerStart - функция запуска HTTP сервера
func HTTPServerStart(log *logrus.Logger, cfg *config.Config, router chi.Router) {
	log.Infof("starting %s", cfg.Server.Address)
	log.Fatal(http.ListenAndServe(cfg.Server.Address, router))
}
