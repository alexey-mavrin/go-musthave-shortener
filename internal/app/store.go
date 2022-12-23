package app

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
)

type store struct {
	s storage
}

type storage interface {
	store(value string) (key string, e error)
	get(key string) (value string, e error)
	open(back string)
	error() error
	close()
}

type mapStorage struct {
	s       map[string]string
	back    *os.File
	mu      sync.Mutex
	lastErr error
}

const (
	keyLen          = 6
	maxStoreAttempt = 10
)

func newStoreWithFile(path string) store {
	stor := store{
		s: storage(&mapStorage{
			s: make(map[string]string),
		}),
	}
	if path != "" {
		stor.s.open(path)
	}
	return stor
}

func (s *mapStorage) error() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	err := s.lastErr
	s.lastErr = nil
	return err
}

func (s *mapStorage) open(path string) {
	log.Printf("opening %s as backing store", path)
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		log.Printf("cannon open file %s: %v", path, err)
		s.mu.Lock()
		s.lastErr = err
		s.back = nil
		s.mu.Unlock()
		return
	}
	s.mu.Lock()
	s.lastErr = nil
	s.back = file
	s.mu.Unlock()
	s.loadFromFile()
}

func (s *mapStorage) close() {
	if s.back != nil {
		s.back.Close()
	}
}

func (s *mapStorage) loadFromFile() {
	if s.back == nil {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.back.Seek(0, 0)
	scanner := bufio.NewScanner(s.back)

	for scanner.Scan() {
		line := scanner.Text()
		toks := strings.Split(line, " ")
		if len(toks) != 2 {
			s.lastErr = fmt.Errorf("cannot split line '%s'", line)
			return
		}
		s.s[toks[0]] = toks[1]
	}

}

func (s *mapStorage) store(url string) (string, error) {
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

	if s.back != nil {
		line := key + " " + url + "\n"
		_, err := s.back.Write([]byte(line))
		if err != nil {
			s.lastErr = err
		}
	}
	return key, nil
}

func (s *mapStorage) get(key string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	url, ok := s.s[key]
	if !ok {
		return "", errors.New("no key exists")
	}
	return url, nil
}
