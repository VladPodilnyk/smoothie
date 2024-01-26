package storage

import (
	"context"
	"errors"
	"fmt"
)

// TODO: fixme (doesn't count windows properly)
type inMemoryStorage struct {
	data map[string]uint
}

func NewInMemoryStorage() inMemoryStorage {
	return inMemoryStorage{data: make(map[string]uint)}
}

func (s inMemoryStorage) Inc(ctx context.Context, key string) {
	_, ok := s.data[key]
	if !ok {
		s.data[key] = 0
	}
	s.data[key]++
}

func (s inMemoryStorage) Get(ctx context.Context, key string) (uint, error) {
	data, ok := s.data[key]
	if !ok {
		return 0, errors.New(fmt.Sprintf("Key %s not found", key))
	}
	return data, nil
}
