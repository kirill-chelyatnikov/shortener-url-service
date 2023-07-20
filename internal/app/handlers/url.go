package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/jackc/pgx/v5"
	"github.com/kirill-chelyatnikov/shortener-url-service/internal/app/models"
	"github.com/kirill-chelyatnikov/shortener-url-service/pkg"
	"io"
	"net/http"
)

// postHandler - функция-хэндлер для обработки POST запросов, отслеживаемый путь: "/"
func (h *Handler) postHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	cookieHash, err := r.Cookie("hash")

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		h.log.Errorf("can't get hash from cookie, err: %s", err)

		return
	}
	//читаем тело запроса
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		h.log.Errorf("unable to read request body, err: %s", err)

		return
	}

	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			h.log.Errorf("can't close Body, err: %s", err)
		}
	}(r.Body)

	//проверка на пустоту тела запроса
	if len(body) == 0 {
		http.Error(w, "empty request body", http.StatusBadRequest)
		h.log.Error("empty request body")

		return
	}

	//если тело запроса прочитано успешно, то генерируем ссылку и записываем её в хранилище
	link := &models.Link{
		BaseURL: string(body),
		Hash:    cookieHash.Value,
	}

	var updated bool
	updated, err = h.service.Add(ctx, link)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		h.log.Error(err)

		return
	}

	responseStatusCode := http.StatusCreated

	//меняем статус код в зафисимости от того добавили мы запись или проапдейтили
	if updated {
		responseStatusCode = http.StatusConflict
	}
	w.WriteHeader(responseStatusCode)

	//записываем ссылку в тело ответа
	_, err = w.Write([]byte(fmt.Sprintf("%s/%s", h.cfg.App.BaseURL, link.ID)))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		h.log.Errorf("failed to write response body, err: %s", err)

		return
	}
	h.log.Infof("URL updated successfully, set %d response code", responseStatusCode)
}

// getHandler - функция-хэндлер для обработки GET запросов, отслеживаемый путь: "/{id}"
func (h *Handler) getHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	//получаем интидификатор ссылки из GET-параметра
	id := chi.URLParam(r, "id")

	/*
		отсутствие проверки на пустоту передаваемого интидификатора обусловлена тем,
		что роут GET "/" не зарегистрован в приложении. По дефолту отдается ошибка 405.
	*/

	//получение url по индетификатору, проверка на его существование
	url, err := h.service.Get(ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		h.log.Error(err)

		return
	}

	w.Header().Set("Location", url)
	w.WriteHeader(http.StatusTemporaryRedirect)
	h.log.Infof("successful redirect to: %s", url)
}

// apiHandler - функция-хэндлер для обработки POST запросов, отслеживаемый путь: "/api/shorten"
func (h *Handler) apiHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	cookieHash, err := r.Cookie("hash")

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		h.log.Errorf("can't get hash from cookie, err: %s", err)

		return
	}

	//читаем тело запроса
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		h.log.Errorf("unable to read request body, err: %s", err)

		return
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			h.log.Errorf("can't close body request, err: %s", err)
		}
	}(r.Body)

	//проверка на пустоту тела запроса
	if len(body) == 0 {
		http.Error(w, "empty request body", http.StatusBadRequest)
		h.log.Error("empty request body")

		return
	}

	//создаём структуры для получения и отправки данных
	apiHandlerRequest := &APIHandlerRequest{}
	apiHandlerResponse := &APIHandlerResponse{}

	/*
		записываем полученный json-объект в заранее созданную структуру.
		Если на вход будет принят неккоректный ключ (не "url"), то ошибка возникнет на моменте добавления урла в хранилище,
		т.к. значение в apiHandlerRequest.URL будет пустое
	*/
	if err = json.Unmarshal(body, apiHandlerRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		h.log.Errorf("cant't unmarshal request body, err: %s", err)

		return
	}

	//генерируем ссылку и записываем её в хранилище
	link := &models.Link{
		BaseURL: apiHandlerRequest.URL,
		Hash:    cookieHash.Value,
	}

	var updated bool
	updated, err = h.service.Add(ctx, link)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		h.log.Errorf("cant't add URL, err: %s", err)

		return
	}

	//записываем результат в структуру ответа
	apiHandlerResponse.Result = fmt.Sprintf("%s/%s", h.cfg.App.BaseURL, link.ID)

	//записываем результат в виде json-объекта
	result, err := json.Marshal(apiHandlerResponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		h.log.Errorf("cant't marshal result, err: %s", err)

		return
	}

	//устанавливаем заголовок "application/json" и код ответа
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	responseStatusCode := http.StatusCreated

	//меняем статус код в зафисимости от того добавили мы запись или проапдейтили
	if updated {
		responseStatusCode = http.StatusConflict
	}
	w.WriteHeader(responseStatusCode)

	//записываем результат в тело ответа
	_, err = w.Write(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		h.log.Errorf("failed to write response body, err: %s", err)

		return
	}

}

func (h *Handler) apiGetAllURLS(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	cookieHash, err := r.Cookie("hash")

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		h.log.Errorf("can't get hash from cookie, err: %s", err)

		return
	}

	links, err := h.service.GetAll(ctx, cookieHash.Value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNoContent)
		h.log.Info(err.Error())

		return
	}

	var result []*APIGETAllResponse

	for _, v := range links {
		result = append(result, &APIGETAllResponse{
			ShortURL:    fmt.Sprintf("%s/%s", h.cfg.App.BaseURL, v.ID),
			OriginalURL: v.BaseURL,
		})
	}

	jsonResult, err := json.Marshal(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		h.log.Errorf("cant't marshal result, err: %s", err)

		return
	}

	//устанавливаем заголовок "application/json" и код ответа
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	//записываем результат в тело ответа
	_, err = w.Write(jsonResult)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		h.log.Errorf("failed to write response body, err: %s", err)

		return
	}
}

func (h *Handler) pingDB(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	conn, err := pgx.Connect(ctx, h.cfg.DB.CDN)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		h.log.Errorf("unable to connect to database, err: %s", err)
		return
	}

	defer func() {
		err = conn.Close(ctx)
		if err != nil {
			h.log.Errorf("can't close database connection, err: %s", err)
		}
	}()

	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte("database is working!"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		h.log.Errorf("failed to write response body, err: %s", err)

		return
	}
	h.log.Info("successful database ping")
}

func (h *Handler) apiBatch(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	cookieHash, err := r.Cookie("hash")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		h.log.Errorf("can't get hash from cookie, err: %s", err)

		return
	}

	//читаем тело запроса
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		h.log.Errorf("unable to read request body, err: %s", err)

		return
	}

	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			h.log.Errorf("can't close body request, err: %s", err)
		}
	}(r.Body)

	//проверка на пустоту тела запроса
	if len(body) == 0 {
		http.Error(w, "empty request body", http.StatusBadRequest)
		h.log.Error("empty request body")

		return
	}

	requestLinks := make([]APIBatchRequest, 0)
	err = json.Unmarshal(body, &requestLinks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		h.log.Errorf("cant't unmarshal request, err: %s", err)

		return
	}

	links := make([]*models.Link, 0)
	for _, v := range requestLinks {
		if len(v.CorrelationID) > 0 && len(v.OriginalURL) > 0 {
			link := &models.Link{
				ID:            pkg.GenerateRandomString(),
				BaseURL:       v.OriginalURL,
				CorrelationID: v.CorrelationID,
				Hash:          cookieHash.Value,
			}

			links = append(links, link)
		}
	}

	err = h.service.AddBatch(ctx, links)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		h.log.Errorf("can't add links batch, err: %s", err)

		return
	}

	result := make([]APIBatchResponse, 0)
	for _, v := range links {
		result = append(result, APIBatchResponse{
			CorrelationID: v.CorrelationID,
			ShortURL:      fmt.Sprintf("%s/%s", h.cfg.App.BaseURL, v.ID),
		})
	}

	jsonResult, err := json.Marshal(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		h.log.Errorf("cant't marshal response, err: %s", err)

		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)

	_, err = w.Write(jsonResult)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		h.log.Errorf("failed to write response body, err: %s", err)

		return
	}
}
