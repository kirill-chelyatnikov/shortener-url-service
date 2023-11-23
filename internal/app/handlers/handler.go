package handlers

import (
	"compress/gzip"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"go.uber.org/zap"
	"io"
	"net/http"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/kirill-chelyatnikov/shortener-url-service/internal/app/models"
	"github.com/kirill-chelyatnikov/shortener-url-service/internal/config"
	"github.com/kirill-chelyatnikov/shortener-url-service/pkg"
	"golang.org/x/exp/slices"
)

const (
	HomeURL    = "/"
	DecodeURL  = "/{id}"
	APIURL     = "/api/shorten"
	APIALLURLS = "/api/user/urls"
	PING       = "/ping"
	APIBATCH   = "/api/shorten/batch"
)

var CookieKey = []byte("cookie_key_7385746739")

var ContentTypesToEncode = []string{
	"application/javascript",
	"application/json",
	"text/css",
	"text/html",
	"text/plain",
	"text/xml",
}

type Handler struct {
	log     *zap.SugaredLogger
	cfg     *config.Config
	service serviceInterface
}

type APIHandlerRequest struct {
	URL string `json:"url"`
}

type APIHandlerResponse struct {
	Result string `json:"result"`
}

type APIGETAllResponse struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type APIBatchRequest struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

type APIBatchResponse struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

type serviceInterface interface {
	Add(ctx context.Context, link *models.Link) (bool, error)
	AddBatch(ctx context.Context, links []*models.Link) error
	Get(ctx context.Context, id string) (string, error)
	GetAll(ctx context.Context, hash string) ([]*models.Link, error)
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

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authCookie, _ := r.Cookie("auth")
		hashCookie, _ := r.Cookie("hash")
		if authCookie == nil || hashCookie == nil || !verifyCookie(authCookie.Value, hashCookie.Value) {
			setAuthCookie(w, r)
		}

		next.ServeHTTP(w, r)
	})
}

func setAuthCookie(w http.ResponseWriter, r *http.Request) *http.Cookie {
	cookie := pkg.GenerateRandomString()
	authCookie := &http.Cookie{Name: "auth", Value: cookie, Path: "/"}
	hashCookie := &http.Cookie{Name: "hash", Value: encryptCookie([]byte(cookie)), Path: "/"}
	http.SetCookie(w, authCookie)
	http.SetCookie(w, hashCookie)
	r.Header.Set("Cookie", fmt.Sprintf("auth=%s; hash=%s", authCookie.Value, hashCookie.Value))
	return hashCookie
}

func encryptCookie(cookie []byte) string {
	hmacCookie := hmac.New(sha256.New, CookieKey)
	hmacCookie.Write(cookie)

	return hex.EncodeToString(hmacCookie.Sum(nil))
}

func verifyCookie(authCookie, hashCookie string) bool {
	hashCookieBytes, _ := hex.DecodeString(hashCookie)

	hm := hmac.New(sha256.New, CookieKey)
	hm.Write([]byte(authCookie))

	res := hmac.Equal(hm.Sum(nil), hashCookieBytes)

	return res
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
	router.Use(AuthMiddleware)

	router.Post(HomeURL, h.postHandler)
	router.Post(APIURL, h.apiHandler)
	router.Post(APIBATCH, h.apiBatch)
	router.Get(DecodeURL, h.getHandler)
	router.Get(APIALLURLS, h.apiGetAllURLS)
	router.Get(PING, h.pingDB)

	return router
}

func NewHandler(log *zap.SugaredLogger, cfg *config.Config, service serviceInterface) *Handler {
	return &Handler{
		log:     log,
		cfg:     cfg,
		service: service,
	}
}
