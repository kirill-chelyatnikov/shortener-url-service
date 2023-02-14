package handlers

import (
	"fmt"
	"github.com/go-chi/chi"
	"io"
	"net/http"
)

// postHandler - функция-хэндлер для обработки POST запросов, отслеживаемый путь: "/"
func (h *Handler) postHandler(w http.ResponseWriter, r *http.Request) {
	//читаем тело запроса
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), 400)
		h.log.Errorf("unable to read request body, err: %s", err)

		return
	}
	defer r.Body.Close()

	//проверка на пустоту тела запроса
	if len(body) == 0 {
		http.Error(w, "empty request body", 400)
		h.log.Error("empty request body")

		return
	}

	//если тело запроса прочитано успешно, то генерируем ссылку и записываем её в хранилище
	generatedURL := h.service.GenerateShortURL()
	h.service.Add(generatedURL, string(body))

	//устанавливаем статус-код 201
	w.WriteHeader(http.StatusCreated)

	//записываем ссылку в тело ответа
	_, err = w.Write([]byte(fmt.Sprintf("http://%s:%d/%s", h.cfg.Server.Address, h.cfg.Server.Port, generatedURL)))
	if err != nil {
		http.Error(w, err.Error(), 500)
		h.log.Errorf("failed to write response body, err: %s", err)

		return
	}
	h.log.Info("URL generated successfully, written to response body")
}

// getHandler - функция-хэндлер для обработки GET запросов, отслеживаемый путь: "/{id}"
func (h *Handler) getHandler(w http.ResponseWriter, r *http.Request) {
	//получаем интидификатор ссылки из GET-параметра
	id := chi.URLParam(r, "id")

	/*
		отсутствие проверки на пустоту передаваемого интидификатора обусловлена тем,
		что роут GET "/" не зарегистрован в приложении. По дефолту отдается ошибка 405.
	*/

	//получение url по индетификатору, проверка на его существование
	url, err := h.service.Get(id)
	if err != nil {
		http.Error(w, err.Error(), 400)
		h.log.Error(err)
		return
	}

	w.Header().Set("Location", url)
	w.WriteHeader(http.StatusTemporaryRedirect)
	h.log.Infof("successful redirect to: %s", url)
}
