package handlers

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/kirill-chelyatnikov/shortener-url-service/internal/app/models"
	"github.com/kirill-chelyatnikov/shortener-url-service/internal/config"
	"github.com/sirupsen/logrus"
)

const (
	HomeURL   = "/"
	DecodeURL = "/{id}"
	ApiURL    = "/api/shorten"
)

type Handler struct {
	log     *logrus.Logger
	cfg     *config.Config
	service ServiceInterface
}

type ApiHandlerRequest struct {
	URL string `json:"url"`
}

type ApiHandlerResponse struct {
	Result string `json:"result"`
}

type ServiceInterface interface {
	Add(link *models.Link) error
	Get(id string) (string, error)
	GenerateShortURL() string
}

func (h *Handler) InitRoutes() chi.Router {
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Post(HomeURL, h.postHandler)
	router.Post(ApiURL, h.apiHandler)
	router.Get(DecodeURL, h.getHandler)

	return router
}

func NewHandler(log *logrus.Logger, cfg *config.Config, service ServiceInterface) *Handler {
	return &Handler{
		log:     log,
		cfg:     cfg,
		service: service,
	}
}
