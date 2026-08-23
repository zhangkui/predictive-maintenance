package service

import (
	"sync"
	"time"
)

type CacheItem struct {
	Value     any
	ExpiresAt time.Time
}
type CacheService struct {
	mu    sync.RWMutex
	items map[string]CacheItem
	ttl   time.Duration
}

func NewCache(ttl time.Duration) *CacheService {
	if ttl <= 0 {
		ttl = 5 * time.Minute
	}
	return &CacheService{items: map[string]CacheItem{}, ttl: ttl}
}
func (c *CacheService) Get(key string) (any, bool) {
	c.mu.RLock()
	item, ok := c.items[key]
	c.mu.RUnlock()
	if !ok || time.Now().After(item.ExpiresAt) {
		if ok {
			c.Delete(key)
		}
		return nil, false
	}
	return item.Value, true
}
func (c *CacheService) Set(key string, value any) {
	c.mu.Lock()
	c.items[key] = CacheItem{Value: value, ExpiresAt: time.Now().Add(c.ttl)}
	c.mu.Unlock()
}
func (c *CacheService) Delete(key string) { c.mu.Lock(); delete(c.items, key); c.mu.Unlock() }
func (c *CacheService) Clear()            { c.mu.Lock(); c.items = map[string]CacheItem{}; c.mu.Unlock() }
func (c *CacheService) Size() int         { c.mu.RLock(); defer c.mu.RUnlock(); return len(c.items) }
func (c *CacheService) Remember(key string, load func() (any, error)) (any, error) {
	if value, ok := c.Get(key); ok {
		return value, nil
	}
	value, e := load()
	if e != nil {
		return nil, e
	}
	c.Set(key, value)
	return value, nil
}
