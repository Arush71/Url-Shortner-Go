package shortner

import (
	"math/rand"
)

func (s *Storage) ShortenUrl() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	for {
		b := make([]byte, 6)
		for i := range 6 {
			randN := rand.Intn(len(charset))
			b[i] = charset[randN]
		}
		strB := string(b)
		if s.Map == nil {
			return strB
		}
		if _, exist := s.Map[strB]; exist {
			continue
		}
		return strB

	}
}
