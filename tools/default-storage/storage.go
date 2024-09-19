// Package defaultstorage provides a thread-safe key-value storage.
package defaultstorage

import (
	"sync"
)

// Storage is a generic thread-safe storage for storing key-value pairs.
// It uses a map to store data and a mutex to ensure synchronization during access.
type Storage struct {
	data map[uint64]any
	mu   *sync.Mutex
}

// New creates and returns a new instance of Storage.
// The storage is initialized with a map and a mutex to ensure thread-safe access.
func New[T any]() *Storage {
	return &Storage{
		data: make(map[uint64]any),
		mu:   &sync.Mutex{},
	}
}

// Set inserts or updates the value associated with the given key.
// It locks the storage to ensure that only one goroutine can modify the data at a time.
func (s *Storage) Set(key uint64, val any) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data[key] = val
}

// Get retrieves the value associated with the given key.
// It locks the storage for thread-safe access and returns the value and a boolean
// indicating whether the key exists in the storage.
func (s *Storage) Get(key uint64) (val any, exists bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	val, exists = s.data[key]
	return
}

// GetAll returns a slice containing all values stored in the storage.
// It locks the storage to ensure thread-safe access while reading the data.
func (s *Storage) GetAll() []any {
	s.mu.Lock()
	defer s.mu.Unlock()

	res := make([]any, 0, len(s.data))

	for _, v := range s.data {
		res = append(res, v)
	}

	return res
}
