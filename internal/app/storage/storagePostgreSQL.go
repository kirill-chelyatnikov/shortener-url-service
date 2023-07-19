package storage

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kirill-chelyatnikov/shortener-url-service/internal/app/models"
	"github.com/kirill-chelyatnikov/shortener-url-service/internal/config"
	"github.com/sirupsen/logrus"
)

type PostgreSQLStorage struct {
	log  *logrus.Logger
	cfg  *config.Config
	pool *pgxpool.Pool
}

// AddURL - функция записи данных в storage (PostgreSQL)
func (p *PostgreSQLStorage) AddURL(ctx context.Context, link *models.Link) error {
	qInsertLink := `
		INSERT INTO links as ls (id, baseurl, hash)
		VALUES
		($1, $2, ARRAY[$3])
		`
	_, err := p.pool.Exec(ctx, qInsertLink, link.ID, link.BaseURL, link.Hash)
	if err != nil {
		return NewDBError("AddURL", "can't do query", err)
	}

	return nil
}

// GetURLByID - функция получения записи из storage (PostgreSQL)
func (p *PostgreSQLStorage) GetURLByID(ctx context.Context, id string) (string, error) {
	var res string
	q := `
	SELECT 
	    baseurl 
	FROM links
	WHERE 
	    id = $1
	`

	row := p.pool.QueryRow(ctx, q, id)

	err := row.Scan(&res)
	if err != nil {
		return "", NewDBError("GetURLByID", "can't scan", err)
	}

	return res, nil
}

// GetAllURLSByHash - функция получения всех записей по хешу из storage (PostgreSQL)
func (p *PostgreSQLStorage) GetAllURLSByHash(ctx context.Context, hash string) ([]*models.Link, error) {
	var links []*models.Link
	q := `
	SELECT 
		id, baseurl
	FROM links
	WHERE 
	    hash = $1
	`

	rows, err := p.pool.Query(ctx, q, hash)
	if err != nil {
		return nil, NewDBError("GetAllURLSByHash", "can't do query", err)
	}

	defer rows.Close()

	for rows.Next() {
		var link models.Link
		err = rows.Scan(&link.ID, &link.BaseURL)
		if err != nil {
			return nil, NewDBError("GetAllURLSByHash", "can't scan", err)
		}
		links = append(links, &link)
	}

	return links, nil
}

// AddURLSBatch - функция добавления записей "пачкой"
func (p *PostgreSQLStorage) AddURLSBatch(ctx context.Context, links []*models.Link) error {

	// Начало транзакции
	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return NewDBError("AddURLSBatch", "can't begin tx", err)
	}

	// Обязательный откат транзакции при возникновении ошибок
	defer tx.Rollback(ctx)

	q := `
		INSERT INTO links
		   (id, baseurl, hash)
		VALUES
			($1, $2, $3)
	`

	for _, v := range links {
		_, err = tx.Exec(ctx, q, v.ID, v.BaseURL, v.Hash)
		if err != nil {
			return NewDBError("AddURLSBatch", "can't exec tx", err)
		}
	}

	if err = tx.Commit(ctx); err != nil {
		return NewDBError("AddURLSBatch", "can't commit tx", err)
	}

	return nil
}

func (p *PostgreSQLStorage) Close() error {
	p.pool.Close()
	return nil
}

// CheckBaseURLExist - функция для проверки нахождения URL в БД
func (p *PostgreSQLStorage) CheckBaseURLExist(ctx context.Context, link *models.Link) (bool, error) {
	q := `
	SELECT id FROM links 
	WHERE baseurl = $1
	`
	var id string

	row := p.pool.QueryRow(ctx, q, link.BaseURL)

	switch err := row.Scan(&id); err {
	// Если ошибка ErrNoRows, значит запись отсутствует, не отдаем ошибку и возвращаем false
	case pgx.ErrNoRows:
		return false, nil
	//если ошибки нет, то запись уже присуствует в БД, возвращаем true присваиваем объекту ID из БД
	case nil:
		link.ID = id
		return true, nil
	default:
		return false, NewDBError("CheckBaseURLExist", "can't scan", err)
	}
}

// UpdateHash - функция для добавления хеша пользователя в уже существующую запись
func (p *PostgreSQLStorage) UpdateHash(ctx context.Context, link *models.Link) error {
	qUpdateHash := `
		UPDATE links ls SET 
		hash = array_append(ls.hash, $1)
		WHERE
		     ls.baseurl = $2
		AND NOT
    	$1 = ANY (ls.hash)
		`
	_, err := p.pool.Exec(ctx, qUpdateHash, link.Hash, link.BaseURL)
	if err != nil {
		return NewDBError("UpdateHash", "can't do query", err)
	}

	return nil
}

// dbConnect - фукция подключения к БД
func dbConnect(ctx context.Context, cfg *config.Config) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, cfg.Db.CDN)
	if err != nil {
		return nil, NewDBError("dbConnect", "DB connect failed", err)
	}

	return pool, nil
}

func NewPostgreSQLStorage(ctx context.Context, log *logrus.Logger, cfg *config.Config) *PostgreSQLStorage {
	pool, err := dbConnect(ctx, cfg)
	if err != nil {
		log.Fatalf("can't create PostgreSQLStorage, err: %s", err)
	}

	//Костыль на случай отсутствия таблицы в БД
	q := `
	CREATE TABLE IF NOT EXISTS links (
		id varchar(10) PRIMARY KEY NOT NULL UNIQUE,
		baseURL text NOT NULL UNIQUE,
		hash varchar(64)[] NOT NULL
	)
	`

	_, err = pool.Exec(ctx, q)
	if err != nil {
		log.Errorf("can't do query, err: %s", err)
	}

	return &PostgreSQLStorage{
		log:  log,
		cfg:  cfg,
		pool: pool,
	}
}
