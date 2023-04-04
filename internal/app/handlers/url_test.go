package handlers

import (
	"fmt"
	"github.com/go-chi/chi"
	"github.com/kirill-chelyatnikov/shortener-url-service/internal/app/models"
	"github.com/kirill-chelyatnikov/shortener-url-service/internal/app/services"
	"github.com/kirill-chelyatnikov/shortener-url-service/internal/app/storage"
	"github.com/kirill-chelyatnikov/shortener-url-service/internal/config"
	"github.com/kirill-chelyatnikov/shortener-url-service/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const testConfigURL = "../../config/config.yml"

//инициализация необходимых зависимостей
var log = logger.InitLogger()
var fl = config.GetFlags()
var cfg = config.GetConfig(log, testConfigURL, fl)
var repository = storage.NewStorage(log, cfg)
var serviceURL = services.NewServiceURL(log, cfg, repository)
var h = NewHandler(log, cfg, serviceURL)

func TestPostHandler(t *testing.T) {
	type want struct {
		code         int
		responseBody string
	}

	tests := []struct {
		name        string
		requestBody string
		want        want
		wantErr     bool
	}{
		{
			name:        "positive test",
			requestBody: "https://testlink.com",
			want: want{
				code: 201,
			},
			wantErr: false,
		},
		{
			name:        "empty body",
			requestBody: "",
			want: want{
				code:         400,
				responseBody: "empty request body",
			},
			wantErr: true,
		},
	}

	router := chi.NewRouter()
	router.Post("/", h.postHandler)
	ts := httptest.NewServer(router)
	defer ts.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request, err := http.NewRequest(http.MethodPost, ts.URL, strings.NewReader(tt.requestBody))
			require.NoError(t, err)

			resp, err := http.DefaultClient.Do(request)
			require.NoError(t, err)

			var body []byte
			body, err = io.ReadAll(resp.Body)
			require.NoError(t, err)
			defer func(Body io.ReadCloser) {
				err = Body.Close()
				require.NoError(t, err)
			}(resp.Body)

			assert.Equal(t, tt.want.code, resp.StatusCode)

			if !tt.wantErr {
				id := strings.TrimPrefix(string(body),
					fmt.Sprintf("http://%s/", cfg.Server.Address))
				_, err = serviceURL.Get(id)
				assert.NoError(t, err)

				return
			}

			assert.Equal(t, tt.want.responseBody, strings.TrimRight(string(body), "\n"))

		})
	}
}

func TestGetHandler(t *testing.T) {
	type want struct {
		code           int
		err            string
		headerLocation string
	}

	tests := []struct {
		name    string
		want    want
		id      string
		wantErr bool
	}{
		{
			name: "wrong id",
			want: want{
				code: 400,
				err:  "can't find URL by id: 123",
			},
			id:      "123",
			wantErr: true,
		},
		{
			name: "positive test",
			want: want{
				code:           307,
				headerLocation: "https://google.com",
			},
			id:      "testUser",
			wantErr: false,
		},
	}

	router := chi.NewRouter()
	router.Get("/{id}", h.getHandler)

	err := repository.AddURL(&models.Link{
		ID:      "testUser",
		BaseURL: "https://google.com",
	})
	assert.NoError(t, err)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			request, err := http.NewRequest(http.MethodGet,
				fmt.Sprintf("http://%s:%d/%s", cfg.Server.Address, cfg.Server.Port, tt.id), nil)
			require.NoError(t, err)

			router.ServeHTTP(w, request)

			assert.Equal(t, tt.want.code, w.Code)
			if !tt.wantErr {
				assert.Equal(t, tt.want.headerLocation, w.Header().Get("Location"))

				return
			}

			body, err := io.ReadAll(w.Body)
			require.NoError(t, err)

			assert.Equal(t, tt.want.err, strings.TrimRight(string(body), "\n"))
		})
	}
}

func TestApiHandler(t *testing.T) {
	type want struct {
		code        int
		contentType string
	}

	tests := []struct {
		name    string
		body    string
		want    want
		textErr string
		wantErr bool
	}{
		{
			name: "positive_test",
			body: "{\"url\":\"https://www.google.com\"}",
			want: want{
				code:        201,
				contentType: "application/json; charset=utf-8",
			},
			wantErr: false,
		},
		{
			name: "empty_body",
			body: "",
			want: want{
				code:        400,
				contentType: "text/plain; charset=utf-8",
			},
			textErr: "empty request body",
			wantErr: true,
		},
		{
			name: "incorrect_url_param",
			body: "{\"url123\":\"https://www.google.com\"}",
			want: want{
				code:        400,
				contentType: "text/plain; charset=utf-8",
			},
			textErr: "empty url received",
			wantErr: true,
		},
	}

	router := chi.NewRouter()
	router.Post("/api/shorten", h.apiHandler)
	ts := httptest.NewServer(router)
	defer ts.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/api/shorten", ts.URL), strings.NewReader(tt.body))
			require.NoError(t, err)

			response, err := http.DefaultClient.Do(request)
			require.NoError(t, err)

			assert.Equal(t, tt.want.code, response.StatusCode)
			assert.Equal(t, tt.want.contentType, response.Header.Get("Content-Type"))

			if tt.wantErr {
				var body []byte
				body, err = io.ReadAll(response.Body)
				require.NoError(t, err)
				defer func(Body io.ReadCloser) {
					err = Body.Close()
					require.NoError(t, err)
				}(response.Body)

				assert.Equal(t, tt.textErr, strings.TrimRight(string(body), "\n"))
			}
		})
	}
}
