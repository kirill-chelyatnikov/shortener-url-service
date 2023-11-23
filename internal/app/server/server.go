package server

import (
	"go.uber.org/zap"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/kirill-chelyatnikov/shortener-url-service/internal/config"
)

// HTTPServerStart - функция запуска HTTP сервера
func HTTPServerStart(log *zap.SugaredLogger, cfg *config.Config, router chi.Router) {
	log.Infof("starting %s", cfg.Server.Address)
	log.Fatal(http.ListenAndServe(cfg.Server.Address, router))
}
