package shortner

type Mapping map[string]string

type CounterMap map[string]int
type Storage struct {
	Map     Mapping
	Counter CounterMap
}

func (s *Storage) StoreUrl(url string, destination string) {
	if s.Map == nil {
		s.Map = make(Mapping)
	}
	s.Map[url] = destination
	if s.Counter == nil {
		s.Counter = make(CounterMap)
	}
	if _, exist := s.Counter[url]; !exist {
		s.Counter[url] = 0
	}
}

func (s *Storage) UpdateCounter(url string) {
	if s.Counter == nil {
		s.Counter = make(CounterMap)
	}
	if _, exist := s.Counter[url]; !exist {
		return
	}
	s.Counter[url]++
}
