package handlers

import (
	"encoding/json"
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
	_, err = w.Write([]byte(fmt.Sprintf("http://%s:%d/%s", h.cfg.App.BaseURL, h.cfg.Server.Port, generatedURL)))
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

// apiHandler - функция-хэндлер для обработки POST запросов, отслеживаемый путь: "/api/shorten"
func (h *Handler) apiHandler(w http.ResponseWriter, r *http.Request) {
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

	//создаём структуры для получения и отправки данных
	apiHandlerRequest := &ApiHandlerRequest{}
	apiHandlerResponse := &ApiHandlerResponse{}

	/*
		записываем полученный json-объект в заранее созданную структуру.
		Если на вход будет принят неккоректный ключ (не "url"), то ошибка возникнет на моменте добавления урла в хранилище,
		т.к. значение в apiHandlerRequest.URL будет пустое
	*/
	if err = json.Unmarshal(body, apiHandlerRequest); err != nil {
		http.Error(w, err.Error(), 400)
		h.log.Errorf("cant't unmarshal request body, err: %s", err)

		return
	}

	//генерируем ссылку и записываем её в хранилище
	generatedURL := h.service.GenerateShortURL()
	err = h.service.Add(generatedURL, apiHandlerRequest.URL)
	if err != nil {
		http.Error(w, err.Error(), 400)
		h.log.Errorf("cant't add URL, err: %s", err)

		return
	}

	//записываем результат в структуру ответа
	apiHandlerResponse.Result = fmt.Sprintf("http://%s:%d/%s", h.cfg.App.BaseURL, h.cfg.Server.Port, generatedURL)

	//записываем результат в виде json-объекта
	result, err := json.Marshal(apiHandlerResponse)
	if err != nil {
		http.Error(w, err.Error(), 400)
		h.log.Errorf("cant't marshal result, err: %s", err)

		return
	}

	//устанавливаем заголовок "application/json" и код ответа
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(201)

	//записываем результат в тело ответа
	_, err = w.Write(result)
	if err != nil {
		http.Error(w, err.Error(), 500)
		h.log.Errorf("failed to write response body, err: %s", err)

		return
	}

}
