package database

import (
	"github.com/patrickmn/go-cache"
	"time"
)

type MemoryCache interface {
	GetCache() *cache.Cache            // GetCache returns the cache instance.
	GetExpireTime() time.Duration      // GetExpireTime returns the expiration time for cache entries.
	GetPurgeExpireTime() time.Duration // GetPurgeExpireTime returns the expiration time for purging cache entries.
}

type memoryCache struct {
	dbCache         *cache.Cache
	expireTime      time.Duration
	purgeExpireTime time.Duration
}

// The NewMemoryCache function is a constructor that creates a new instance of memoryCache and initializes its fields
// with the provided expireTime and purgeExpireTime values.
func NewMemoryCache(expireTime time.Duration, purgeExpireTime time.Duration) MemoryCache {
	dbCache := cache.New(
		expireTime,
		purgeExpireTime,
	)
	return &memoryCache{dbCache: dbCache, expireTime: expireTime, purgeExpireTime: purgeExpireTime}
}

// GetCache returns the cache instance.
func (m memoryCache) GetCache() *cache.Cache {
	return m.dbCache
}

// GetExpireTime returns the expiration time for cache entries.
func (m memoryCache) GetExpireTime() time.Duration {
	return m.expireTime
}

// GetPurgeExpireTime returns the expiration time for purging cache entries.
func (m memoryCache) GetPurgeExpireTime() time.Duration {
	return m.purgeExpireTime
}
