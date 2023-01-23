package main

import (
	"github.com/kirill-chelyatnikov/shortener-url-service/internal/app"
	"github.com/kirill-chelyatnikov/shortener-url-service/internal/config"
	"github.com/kirill-chelyatnikov/shortener-url-service/pkg"
)

func main() {
	log := pkg.InitLogger()
	cfg := config.GetConfig(log)
	server := app.NewServer(log, cfg)
	server.HttpServerStart()
}
