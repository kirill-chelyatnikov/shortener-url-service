package handlers

import (
	"compress/gzip"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/kirill-chelyatnikov/shortener-url-service/internal/app/models"
	"github.com/kirill-chelyatnikov/shortener-url-service/internal/config"
	"github.com/sirupsen/logrus"
	"golang.org/x/exp/slices"
	"io"
	"net/http"
	"strings"
)

const (
	HomeURL   = "/"
	DecodeURL = "/{id}"
	ApiURL    = "/api/shorten"
)

var ContentTypesToEncode = []string{
	"application/javascript",
	"application/json",
	"text/css",
	"text/html",
	"text/plain",
	"text/xml",
}

type Handler struct {
	log     *logrus.Logger
	cfg     *config.Config
	service serviceInterface
}

type ApiHandlerRequest struct {
	URL string `json:"url"`
}

type ApiHandlerResponse struct {
	Result string `json:"result"`
}

type serviceInterface interface {
	Add(link *models.Link) error
	Get(id string) (string, error)
	GenerateShortURL() string
}

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (gw gzipWriter) Write(b []byte) (int, error) {
	return gw.Writer.Write(b)
}

type Compressor struct {
	gz *gzip.Writer
}

// GzipMiddlewareResponse - middleware ответа клиенту в формате gzip
func (c *Compressor) GzipMiddlewareResponse(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// проверяем Accept-Encoding и Content-Type
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") ||
			!slices.Contains(ContentTypesToEncode, r.Header.Get("Content-Type")) {

			// если клиент не может принять gzip или передан несоответствующий Content-Type, то пропускаем обработку
			next.ServeHTTP(w, r)
			return
		}

		if c.gz == nil {
			// создаём gzip.Writer поверх текущего w
			gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
			if err != nil {
				_, errWr := io.WriteString(w, err.Error())
				if errWr != nil {
					panic(errWr)
				}
				return
			}
			c.gz = gz
		} else {
			// если Writer уже создан, делаем Reset
			c.gz.Reset(w)
		}

		defer func(gz *gzip.Writer) {
			err := gz.Close()
			if err != nil {
				panic(err)
			}
		}(c.gz)

		w.Header().Set("Content-Encoding", "gzip")
		next.ServeHTTP(&gzipWriter{ResponseWriter: w, Writer: c.gz}, r)
	})
}

// GzipMiddlewareRequest - middleware обработки запроса от клиента в формате gzip
func GzipMiddlewareRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// проверяем Content-Encoding. В случаи успеха - декодированный ответ подствляем в Body
		if r.Header.Get("Content-Encoding") == "gzip" {
			gz, err := gzip.NewReader(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			r.Body = gz
			defer func(gz *gzip.Reader) {
				err = gz.Close()
				if err != nil {
					panic(err)
				}
			}(gz)
		}

		next.ServeHTTP(w, r)
	})
}

func (h *Handler) InitRoutes() chi.Router {
	compressor := &Compressor{}
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(compressor.GzipMiddlewareResponse)
	router.Use(GzipMiddlewareRequest)

	router.Post(HomeURL, h.postHandler)
	router.Post(ApiURL, h.apiHandler)
	router.Get(DecodeURL, h.getHandler)

	return router
}

func NewHandler(log *logrus.Logger, cfg *config.Config, service serviceInterface) *Handler {
	return &Handler{
		log:     log,
		cfg:     cfg,
		service: service,
	}
}
