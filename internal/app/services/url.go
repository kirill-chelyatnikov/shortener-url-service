package services

import (
	"context"
	"errors"
	"time"

	"github.com/kirill-chelyatnikov/shortener-url-service/internal/app/models"
	"github.com/kirill-chelyatnikov/shortener-url-service/pkg"
)

const (
	deleteDelay = 3 * time.Second
	batchCount  = 10
)

// Add - функция сервиса для добавления/изменения записи
func (s *ServiceURL) Add(ctx context.Context, link *models.Link) (bool, error) {
	// Проверка на пустоту переданных полей
	if len(link.BaseURL) == 0 || len(link.Hash) == 0 {
		return false, errors.New("empty url received")
	}

	// Проверяем наличие URL в БД
	check, err := s.repository.CheckBaseURLExist(ctx, link)
	if err != nil {
		return false, err
	}

	/* Если URL в базе уже существует, то добавляем хеш пользователя в массив хэшей данной записи,
	также передаем информацию в хендлер о том, что запись мы проапдейтили, а не вставили новую, нужно для понимания
	какой код ответа следует отдать */
	if check {
		if err = s.repository.UpdateHash(ctx, link); err != nil {
			return false, err
		}
		return true, nil
	}
	/* Если URL не найден, то просто пишем его в БД.
	Присваиваем записи ID (в случае апдейта записи, ID мы берем из БД) */
	link.ID = pkg.GenerateRandomString()
	if err = s.repository.AddURL(ctx, link); err != nil {
		return false, err
	}

	return false, nil
}

// AddBatch - функция сервиса для добавления записей "пачкой"
func (s *ServiceURL) AddBatch(ctx context.Context, links []*models.Link) error {
	if len(links) == 0 {
		return errors.New("passed an empty array of references")
	}

	return s.repository.AddURLSBatch(ctx, links)
}

// Get - функция сервиса для получения записи по ID
func (s *ServiceURL) Get(ctx context.Context, id string) (*models.Link, error) {
	return s.repository.GetURLByID(ctx, id)
}

// GetAll - функция сервиса для получения всех записей по хешу
func (s *ServiceURL) GetAll(ctx context.Context, hash string) ([]*models.Link, error) {
	return s.repository.GetAllURLSByHash(ctx, hash)
}

// DeleteBatch - функция сервиса для удаления записей пачкой
func (s *ServiceURL) DeleteBatch(ctx context.Context, links []string, hash string) error {
	//Проверяем какие URL принадлежат пользователю
	userLinks, err := s.repository.GetAllURLSByHash(ctx, hash)
	if err != nil {
		return err
	}

	idsMap := make(map[string]struct{})
	for _, v := range links {
		idsMap[v] = struct{}{}
	}

	//Проверенные URL пишем в канал
	for _, v := range userLinks {
		if _, ok := idsMap[v.ID]; ok {
			s.deleteCh <- v.ID
		}
	}

	return nil
}

func (s *ServiceURL) CheckBatches(ctx context.Context) {
	idsToDelete := make([]string, 0)
	ticker := time.NewTicker(deleteDelay)
	for {
		select {
		case <-ctx.Done():
			close(s.deleteCh)
			s.log.Info("ch closed")
			return
		//При записи в канал проверяем кол-во значений, если > 10, производим удаление
		case id := <-s.deleteCh:
			idsToDelete = append(idsToDelete, id)
			if len(idsToDelete) < batchCount {
				continue
			}
		//Производим удаление каждые 15 секунд
		case <-ticker.C:
		default:
			continue
		}

		if len(idsToDelete) > 0 {
			err := s.repository.DeleteURLSBatch(ctx, idsToDelete)
			if err != nil {
				s.log.Errorf("can't delete id's, err: %s", err)
				return
			}
			s.log.Info("successful deleted")
			idsToDelete = nil
		} else {
			s.log.Info("nothing to delete")
		}
	}
}
