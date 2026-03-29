package cache

import (
	"context"
	"maps"
	"sync"
	"time"

	"github.com/Arush71/url-shortener/internal/db"
)

type CacheMap map[string]string
type CounterMap map[string]int64

type Cache struct {
	CMap     CacheMap
	CounterM CounterMap
	mu       sync.RWMutex
	Q        *db.Queries
}

func (c *Cache) GetUrl(code string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
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
	c.CounterM[code]++
}

func (c *Cache) Flush() {
	c.mu.Lock()
	localCopy := make(map[string]int64)
	maps.Copy(localCopy, c.CounterM)
	c.CounterM = make(CounterMap)
	c.mu.Unlock()
	for k, v := range localCopy {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
		err := c.Q.UpdateCounter(ctx, db.UpdateCounterParams{
			Code:    k,
			Counter: v,
		})
		cancel()
		if err == nil {
			delete(localCopy, k)
		}
	}
	if len(localCopy) != 0 {
		c.mu.Lock()
		defer c.mu.Unlock()
		for k, v := range localCopy {
			c.CounterM[k] += v
		}
	}
}

func (c *Cache) GetCounter(code string) (int64, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	val, ok := c.CounterM[code]
	return val, ok
}
