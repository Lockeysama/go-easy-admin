package utils

import (
	"fmt"
	"time"

	cache "github.com/patrickmn/go-cache"
)

// Cache 全局缓存实例
var _cache map[string]*cache.Cache

// DefaultExpiration 默认过期时间
var DefaultExpiration time.Duration

// DefaultCleanupInterval 默认缓存清除时间
var DefaultCleanupInterval time.Duration

// DefaultCacheName 默认缓存实例名称
const DefaultCacheName = "default"

// InitCache 初始化缓存
func InitCache() {
	_cache = make(map[string]*cache.Cache)
	if DefaultExpiration == 0 {
		DefaultExpiration = time.Hour
	}
	if DefaultCleanupInterval == 0 {
		DefaultCleanupInterval = time.Hour * 2
	}
	NewCache(DefaultCacheName, DefaultExpiration, DefaultCleanupInterval)
}

// NewCache 新建缓存实例
func NewCache(name string, expiration time.Duration, cleanupInterval time.Duration) *cache.Cache {
	if _, ok := _cache[name]; !ok {
		_cache[name] = cache.New(expiration, cleanupInterval)
	}
	return _cache[name]
}

// Cache 根据名称获取 Cache 实例
func Cache(name string) *cache.Cache {
	if _, ok := _cache[name]; !ok {
		panic(fmt.Sprintf("cache \"%s\" not found", name))
	}
	return _cache[name]
}

// DefaultCache 获取默认 Cache 实例
func DefaultCache() *cache.Cache {
	if defaultCache, ok := _cache[DefaultCacheName]; ok {
		return defaultCache
	}
	panic("default cache not found, please called \"InitCache\" method.")
}
