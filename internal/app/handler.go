package app

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"io"
	"math/rand"
	"net/http"
	"time"
)

//хранилище
var data = make(map[string]string)

func (s *server) postHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	//проверка на пустоту тела запроса
	if r.Body == http.NoBody {
		http.Error(w, "empty request body", 400)
		s.log.Error("empty request body")

		return
	}

	//читаем тело запроса
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), 400)
		s.log.Errorf("unable to read request body, err: %s", err)

		return
	}

	//если тело запроса прочитано успешно, то генерируем ссылку и записываем её в хранилище
	generatedURL := s.generateShortURL()
	data[generatedURL] = string(body)

	w.WriteHeader(http.StatusCreated)

	//записываем ссылку в тело ответа
	_, err = w.Write([]byte(fmt.Sprintf("http://%s:%s/%s", s.cfg.Server.Address, s.cfg.Server.Port, generatedURL)))
	if err != nil {
		http.Error(w, err.Error(), 400)
		s.log.Errorf("failed to write response body, err: %s", err)

		return
	}
	s.log.Info("URL generated successfully, written to response body")
}

func (s *server) getHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	//получаем интидификатор ссылки из GET-параметра
	id := p.ByName("id")

	/*
		отсутствие проверки на пустоту передаваемого интидификатора обусловлена тем,
		что роут GET "/" не зарегистрован в приложении. По дефолту отдается ошибка 405.
	*/

	//проверка нахождения интидификатора в хранилище
	if _, ok := data[id]; !ok {
		http.Error(w, "can't find url", 400)
		s.log.Errorf("can't find url in storage. Id: %s", id)

		return
	}

	w.Header().Set("Location", data[id])
	w.WriteHeader(http.StatusTemporaryRedirect)
	s.log.Infof("successful redirect to: %s", data[id])
}

// generateShortURL - функция генерации короткого URL
func (s *server) generateShortURL() string {
	rand.Seed(time.Now().UnixNano())
	var chars = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0987654321")
	str := make([]rune, s.cfg.App.ShortedURLLen)

	for i := range str {
		str[i] = chars[rand.Intn(len(chars))]
	}

	return string(str)
}
