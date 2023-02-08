package server

import (
	"fmt"
	"github.com/go-chi/chi"
	"github.com/kirill-chelyatnikov/shortener-url-service/internal/config"
	"github.com/sirupsen/logrus"
	"net/http"
)

// HTTPServerStart - функция запуска HTTP сервера
func HTTPServerStart(log *logrus.Logger, cfg *config.Config, router chi.Router) {
	log.Infof("starting server on port %d", cfg.Server.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", cfg.Server.Address, cfg.Server.Port), router))
}
