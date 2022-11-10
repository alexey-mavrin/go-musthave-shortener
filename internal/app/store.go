package app

import (
	"errors"
	"sync"
)

type storage struct {
	s  map[string]string
	mu sync.Mutex
}

const (
	keyLen          = 6
	maxStoreAttempt = 10
)

func newStorage() *storage {
	stor := make(map[string]string)
	return &storage{
		s: stor,
	}
}

func (s *storage) store(url string) (string, error) {
	var key string

	s.mu.Lock()
	defer s.mu.Unlock()

	for i := 0; i < maxStoreAttempt; i++ {
		key = randSeq(keyLen)
		if _, ok := s.s[key]; !ok {
			break
		}
		key = ""
	}

	if key == "" {
		return "", errors.New("cannot generate storage key")
	}

	s.s[key] = url
	return key, nil
}

func (s *storage) get(key string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	url, ok := s.s[key]
	if !ok {
		return "", errors.New("no key exists")
	}
	return url, nil
}
