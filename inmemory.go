package mptrie

import (
	"errors"
	"sync"
)

var (
	KeyNotFound    = errors.New("cannot found key")
	NotInitialized = errors.New("database not initialized")
)

type InMemoryStorage struct {
	kv   map[string][]byte
	lock sync.RWMutex
}

func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		kv: make(map[string][]byte),
	}
}

func (s *InMemoryStorage) Put(key, value []byte) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	if s.kv == nil {
		return NotInitialized
	}

	s.kv[string(key)] = value
	return nil
}

func (s *InMemoryStorage) Delete(key []byte) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	if s.kv == nil {
		return NotInitialized
	}

	delete(s.kv, string(key))
	return nil
}

func (s *InMemoryStorage) Has(key []byte) (bool, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	if s.kv == nil {
		return false, NotInitialized
	}

	_, ok := s.kv[string(key)]
	return ok, nil
}

func (s *InMemoryStorage) Get(key []byte) ([]byte, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	if s.kv == nil {
		return nil, NotInitialized
	}

	if entry, ok := s.kv[string(key)]; ok {
		if entry == nil {
			return []byte{}, nil
		}

		newb := make([]byte, len(entry))
		copy(newb[:], entry)

		return newb, nil
	}

	return nil, KeyNotFound
}
