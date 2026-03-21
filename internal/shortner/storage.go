package shortner

import (
	"fmt"
	"sync"
)

type Mapping map[string]string

type CounterMap map[string]int
type Storage struct {
	Map     Mapping
	Counter CounterMap
	Mu      sync.Mutex
}

func (s *Storage) Store(url string, destination string) {
	s.Map[url] = destination
	if _, exist := s.Counter[url]; !exist {
		s.Counter[url] = 0
	}
}

func (s *Storage) UpdateCounter(short_url string) {
	if _, exist := s.Counter[short_url]; !exist {
		return
	}
	s.Counter[short_url]++
}

func (s *Storage) GetDestination(short_url string) *string {
	if value, ok := s.Map[short_url]; ok {
		return &value
	}
	return nil
}

func CreateStorage() *Storage {
	return &Storage{
		Map:     make(Mapping),
		Counter: make(CounterMap),
	}
}

func (s *Storage) GetStats(short_url string) (string, int, error) {
	if url, ok := s.Map[short_url]; ok {
		if counter, ok := s.Counter[short_url]; ok {
			return url, counter, nil
		}
	}
	return "", 0, fmt.Errorf("Url doesn't exist")
}
