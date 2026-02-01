package store

import (
	"errors"
	"sync"

	"github.com/google/uuid"
)

type Store struct {
	m     sync.Mutex
	store []News
}

func New() *Store {
	return &Store{
		m:     sync.Mutex{},
		store: []News{},
	}
}
func (s *Store) Create(newsRecord News) (News, error) {
	s.m.Lock()
	defer s.m.Unlock()
	id := uuid.New()
	newsRecord.Id = id
	s.store = append(s.store, newsRecord)
	return newsRecord, nil
}

func (s *Store) FindAll() ([]News, error) {
	s.m.Lock()
	defer s.m.Unlock()
	return s.store, nil
}

func (s *Store) FindById(id uuid.UUID) (News, error) {
	s.m.Lock()
	defer s.m.Unlock()
	for _, news := range s.store {
		if news.Id == id {
			return news, nil
		}
	}
	return News{}, errors.New("record not found")
}

func (s *Store) UpdateById(id uuid.UUID, newsRecord News) (News, error) {
	s.m.Lock()
	defer s.m.Unlock()
	for index, news := range s.store {
		if news.Id == id {
			newsRecord.Id = id
			s.store = append(s.store[:index], s.store[index+1:]...)
			s.store = append(s.store, newsRecord)
			return newsRecord, nil
		}
	}
	return News{}, errors.New("update failed, record not found")
}

func (s *Store) DeleteById(id uuid.UUID) error {
	s.m.Lock()
	defer s.m.Unlock()
	for index, news := range s.store {
		if news.Id == id {
			s.store = append(s.store[:index], s.store[index+1:]...)
			return nil
		}
	}
	return errors.New("delete failed, record not found")
}
