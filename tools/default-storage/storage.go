package defaultstorage

import (
	"sync"
)

type Storage struct {
	data map[uint64]any
	mu   *sync.Mutex
}

func New[T any]() *Storage {
	return &Storage{
		data: make(map[uint64]any),
		mu:   &sync.Mutex{},
	}
}

func (s *Storage) Set(key uint64, val any) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data[key] = val
}

func (s *Storage) Get(key uint64) (val any, exists bool) {
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
