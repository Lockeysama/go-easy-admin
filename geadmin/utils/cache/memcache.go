package cache

import (
	"fmt"
	"time"

	"github.com/patrickmn/go-cache"
)

// Cache 全局缓存实例
var _memCache map[string]*cache.Cache

// DefaultMemCacheExpiration 默认过期时间
var DefaultMemCacheExpiration time.Duration

// DefaultCleanupInterval 默认缓存清除时间
var DefaultMemCacheCleanupInterval time.Duration

// DefaultCacheName 默认缓存实例名称
const DefaultMemCacheName = "default"

// InitCache 初始化缓存
func InitCache() {
	_memCache = make(map[string]*cache.Cache)
	if DefaultMemCacheExpiration == 0 {
		DefaultMemCacheExpiration = time.Hour
	}
	if DefaultMemCacheCleanupInterval == 0 {
		DefaultMemCacheCleanupInterval = time.Hour * 2
	}
	NewCache(DefaultMemCacheName, DefaultMemCacheExpiration, DefaultMemCacheCleanupInterval)
}

// NewCache 新建缓存实例
func NewCache(name string, expiration time.Duration, cleanupInterval time.Duration) *cache.Cache {
	if _, ok := _memCache[name]; !ok {
		_memCache[name] = cache.New(expiration, cleanupInterval)
	}
	return _memCache[name]
}

// Cache 根据名称获取 Cache 实例
func MemCache(name string) *cache.Cache {
	if _, ok := _memCache[name]; !ok {
		panic(fmt.Sprintf("cache \"%s\" not found", name))
	}
	return _memCache[name]
}

// DefaultCache 获取默认 Cache 实例
func DefaultMemCache() *cache.Cache {
	if defaultCache, ok := _memCache[DefaultMemCacheName]; ok {
		return defaultCache
	}
	panic("default cache not found, please called \"InitCache\" method.")
}
