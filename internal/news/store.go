package news

import (
	"context"

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
		return createdNews, err
	}
	return createdNews, nil
}

// get all news
func (s Store) FindAll(ctx context.Context) (news []Record, err error) {
	err = s.db.NewSelect().Model(&news).Scan(ctx)
	if err != nil {
		return news, err
	}
	return news, nil
}

// get news by id

func (s Store) FindById(ctx context.Context, id uuid.UUID) (news Record, err error) {
	err = s.db.NewSelect().Model(&news).Where("id = ?", id).Scan(ctx, &news)
	if err != nil {
		return news, err
	}
	return news, nil
}

// update news by id
func (s Store) UpdateById(ctx context.Context, id uuid.UUID, news Record) (err error) {
	_, err = s.db.NewUpdate().Model(&news).Where("id = ?", id).Returning("NULL").Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

// delete news by id

func (s Store) DeleteById(ctx context.Context, id uuid.UUID) (err error) {
	_, err = s.db.NewDelete().Model(&Record{}).Where("id = ?", id).Returning("NULL").Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}
