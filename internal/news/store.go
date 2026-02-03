package news

import (
	"context"
	"database/sql"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Store struct {
	db bun.IDB
}

func NewStore(db bun.IDB) *Store {
	return &Store{
		db: db,
	}
}

// create news
func (s Store) Create(ctx context.Context, news Record) (createdNews Record, err error) {
	news.Id = uuid.New()
	err = s.db.NewInsert().Model(&news).Returning("*").Scan(ctx, &createdNews)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return createdNews, NewCustomError(err, http.StatusNotFound)
		}
		return createdNews, NewCustomError(err, http.StatusInternalServerError)
	}
	return createdNews, nil
}

// get all news
func (s Store) FindAll(ctx context.Context) (news []Record, err error) {
	err = s.db.NewSelect().Model(&news).Scan(ctx)
	if err != nil {
		return news, NewCustomError(err, http.StatusInternalServerError)
	}
	return news, nil
}

// get news by id
func (s Store) FindById(ctx context.Context, id uuid.UUID) (news Record, err error) {
	err = s.db.NewSelect().Model(&news).Where("id = ?", id).Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return news, NewCustomError(err, http.StatusNotFound)
		}
		return news, NewCustomError(err, http.StatusInternalServerError)
	}
	return news, nil
}

// update news by id
func (s Store) UpdateById(ctx context.Context, id uuid.UUID, news Record) error {
	r, err := s.db.NewUpdate().Model(&news).Where("id = ?", id).Returning("NULL").Exec(ctx)
	if err != nil {
		return NewCustomError(err, http.StatusInternalServerError)
	}
	rows, err := r.RowsAffected()
	if err != nil {
		return NewCustomError(err, http.StatusInternalServerError)
	}
	if rows == 0 {
		return NewCustomError(errors.New("record not found"), http.StatusNotFound)
	}
	return nil
}

// delete news by id
func (s Store) DeleteById(ctx context.Context, id uuid.UUID) (err error) {
	_, err = s.db.NewDelete().Model(&Record{}).Where("id = ?", id).Returning("NULL").Exec(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}
		return NewCustomError(err, http.StatusInternalServerError)
	}
	return nil
}
