package handlers

import (
	"fmt"
	"github.com/go-chi/chi"
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
var cfg = config.GetConfig(log, testConfigURL)
var repository = storage.NewStorage()
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

			body, err := io.ReadAll(resp.Body)
			defer resp.Body.Close()
			require.NoError(t, err)

			assert.Equal(t, tt.want.code, resp.StatusCode)

			if !tt.wantErr {
				id := strings.TrimPrefix(string(body),
					fmt.Sprintf("http://%s:%d/", cfg.Server.Address, cfg.Server.Port))
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

	repository.AddURL("testUser", "https://google.com")

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
