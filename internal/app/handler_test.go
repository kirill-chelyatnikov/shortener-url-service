package app

import (
	"github.com/julienschmidt/httprouter"
	config "github.com/kirill-chelyatnikov/shortener-url-service/internal/config"
	"github.com/kirill-chelyatnikov/shortener-url-service/internal/storage"
	"github.com/kirill-chelyatnikov/shortener-url-service/pkg"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

//инициализация необходимых аргументов для тестирования postHandler и getHandler
var log = pkg.InitLogger()
var testCfg = &config.Config{
	Server: struct {
		Address string `yaml:"address"`
		Port    int    `yaml:"port"`
	}{
		Address: "localhost",
		Port:    8080,
	},
	App: struct {
		ShortedURLLen uint8 `yaml:"shortedURLLen"`
	}{
		ShortedURLLen: 5,
	},
}

var rep = storage.NewStorage()
var s = NewServer(log, testCfg, rep)

func Test_server_postHandler(t *testing.T) {
	type want struct {
		code    int
		resBody string
	}
	tests := []struct {
		name    string
		want    want
		reqBody string
		wantErr bool
	}{
		{
			name: "empty body",
			want: want{
				code:    400,
				resBody: "empty request body",
			},
			wantErr: true,
		},
		{
			name: "positive test",
			want: want{
				code: 201,
			},
			reqBody: "https://testlink.com",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.reqBody))

			s.postHandler(w, r, nil)

			body, err := io.ReadAll(w.Body)
			if err != nil {
				t.Log(err)
			}

			if !tt.wantErr {
				//проверка статус-кода при позитивном сценарии
				assert.Equal(t, tt.want.code, w.Code)

				//получение ID строки из body, проверка на существование сгенерированной ссылки по средствам вызова метода GetURLByID
				id := strings.TrimPrefix(string(body), "http://localhost:8080/")
				_, err = s.repository.GetURLByID(id)
				assert.NoError(t, err)

				return
			}

			//проверка статус-вода и body при негативных сценариях
			assert.Equal(t, tt.want.code, w.Code)
			assert.Equal(t, tt.want.resBody, strings.TrimRight(string(body), "\n"))

		})
	}

}

func Test_server_getHandler(t *testing.T) {
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
				headerLocation: "https://testlink.com",
			},
			id:      "testUser",
			wantErr: false,
		},
	}

	rep.AddURL("testUser", "https://testlink.com")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/"+tt.id, nil)

			s.getHandler(w, r, httprouter.Params{httprouter.Param{
				Key:   "id",
				Value: tt.id,
			}})

			_, err := rep.GetURLByID(tt.id)

			assert.Equal(t, tt.want.code, w.Code)
			if !tt.wantErr {
				assert.Equal(t, tt.want.headerLocation, w.Header().Get("Location"))

				return
			}

			assert.Equal(t, tt.want.err, err.Error())
		})
	}

}
