package main

import (
	"flag"
	"github.com/kirill-chelyatnikov/shortener-url-service/internal/app/handlers"
	"github.com/kirill-chelyatnikov/shortener-url-service/internal/app/server"
	"github.com/kirill-chelyatnikov/shortener-url-service/internal/app/services"
	"github.com/kirill-chelyatnikov/shortener-url-service/internal/app/storage"
	"github.com/kirill-chelyatnikov/shortener-url-service/internal/config"
	"github.com/kirill-chelyatnikov/shortener-url-service/pkg/logger"
)

const configURL = "internal/config/config.yml"

func main() {
	fl := config.GetFlags()
	flag.Parse()

	log := logger.InitLogger()
	cfg := config.GetConfig(log, configURL, fl)

	repository := storage.NewStorage(log, cfg)
	defer func() {
		if err := repository.Close(); err != nil {
			log.Errorf("can't close file, err: %s", err)
		}
	}()
	ServiceURL := services.NewServiceURL(log, cfg, repository)
	handler := handlers.NewHandler(log, cfg, ServiceURL)
	server.HTTPServerStart(log, cfg, handler.InitRoutes())
}
