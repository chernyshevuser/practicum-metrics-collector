package storage

import (
	"sync"
)

type Storage struct {
	data map[string]any
	mu   *sync.Mutex
}

func New[T any]() *Storage {
	return &Storage{
		data: make(map[string]any),
		mu:   &sync.Mutex{},
	}
}

func (s *Storage) Set(key string, val any) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data[key] = val
}

func (s *Storage) Get(key string) (val any, exists bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	val, exists = s.data[key]
	return
}

func (s *Storage) GetAll() []any {
	s.mu.Lock()
	defer s.mu.Unlock()

	res := make([]any, 0, len(s.data))

	for _, v := range s.data {
		res = append(res, v)
	}

	return res
}
