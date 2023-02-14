package main

import (
	"github.com/kirill-chelyatnikov/shortener-url-service/internal/app/handlers"
	"github.com/kirill-chelyatnikov/shortener-url-service/internal/app/server"
	"github.com/kirill-chelyatnikov/shortener-url-service/internal/app/services"
	"github.com/kirill-chelyatnikov/shortener-url-service/internal/app/storage"
	"github.com/kirill-chelyatnikov/shortener-url-service/internal/config"
	"github.com/kirill-chelyatnikov/shortener-url-service/pkg/logger"
)

const configURL = "internal/config/config.yml"

func main() {
	log := logger.InitLogger()
	cfg := config.GetConfig(log, configURL)

	repository := storage.NewStorage()
	ServiceURL := services.NewServiceURL(log, cfg, repository)
	handler := handlers.NewHandler(log, cfg, ServiceURL)
	server.HTTPServerStart(log, cfg, handler.InitRoutes())
}
