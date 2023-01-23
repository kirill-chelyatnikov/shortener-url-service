package app

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/kirill-chelyatnikov/shortener-url-service/internal/config"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
)

const (
	HomeURL   = "/"
	DecodeURL = "/:id"
)

type server struct {
	log *logrus.Logger
	cfg *config.Config
}

func (s *server) HttpServerStart() {
	router := httprouter.New()
	router.POST(HomeURL, s.postHandler)
	router.GET(DecodeURL, s.getHandler)

	s.log.Infof("starting server on port %s", s.cfg.Server.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%s", s.cfg.Server.Address, s.cfg.Server.Port), router))
}

func NewServer(log *logrus.Logger, cfg *config.Config) *server {
	return &server{
		log: log,
		cfg: cfg,
	}
}
