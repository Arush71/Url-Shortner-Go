package cache

import (
	"context"
	"sync"
	"time"

	"github.com/Arush71/url-shortener/internal/db"
)

type CacheMap map[string]string
type CounterMap map[string]int64

type Cache struct {
	CMap     CacheMap
	CounterM CounterMap
	mu       sync.Mutex
	Q        *db.Queries
}

func (c *Cache) GetUrl(code string) (string, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if value, ok := c.CMap[code]; ok {
		return value, true
	}
	return "", false
}
func (c *Cache) SaveUrl(code, url string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.CMap[code] = url
}

func SetupCache(Q *db.Queries) *Cache {
	return &Cache{
		CMap:     make(CacheMap),
		CounterM: make(CounterMap),
		Q:        Q,
	}
}
func (c *Cache) IncrementCounter(code string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, ok := c.CounterM[code]; ok {
		c.CounterM[code]++
		return
	}
	c.CounterM[code] = 1
}

func (c *Cache) Flush() {
	c.mu.Lock()
	defer c.mu.Unlock()
	for k, v := range c.CounterM {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
		err := c.Q.UpdateCounter(ctx, db.UpdateCounterParams{
			Code:    k,
			Counter: v,
		})
		cancel()
		if err == nil {
			delete(c.CounterM, k)
		}
	}
}
