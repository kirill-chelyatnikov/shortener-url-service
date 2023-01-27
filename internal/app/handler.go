package app

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"io"
	"math/rand"
	"net/http"
	"time"
)

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
	err = s.repository.AddURL(generatedURL, string(body))
	if err != nil {
		s.log.Error(err)
	}

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

	//получение url по индетификатору, проверка на его существование
	url, err := s.repository.GetURLByID(id)
	if err != nil {
		http.Error(w, "can't find URL", 400)
		s.log.Error(err)
		return
	}

	w.Header().Set("Location", url)
	w.WriteHeader(http.StatusTemporaryRedirect)
	s.log.Infof("successful redirect to: %s", url)
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
