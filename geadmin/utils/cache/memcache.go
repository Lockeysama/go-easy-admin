package cache

import (
	"time"
)

type Cache map[string]interface{}

type Caches map[string]Cache

// Cache 全局缓存实例
var _memCache map[string]Caches

// DefaultMemCacheExpiration 默认过期时间
var DefaultMemCacheExpiration time.Duration

// DefaultCacheName 默认缓存实例名称
const DefaultMemCacheName = "default"

// Cache 根据名称获取 Cache 实例
func MemCache(args ...string) Caches {
	var name string
	if len(args) == 0 {
		name = DefaultMemCacheName
	}
	if _memCache == nil {
		_memCache = make(map[string]Caches)
		if DefaultMemCacheExpiration == 0 {
			DefaultMemCacheExpiration = time.Hour
		}
	}
	if _, ok := _memCache[name]; !ok {
		_memCache[name] = make(Caches)
	}
	return _memCache[name]
}

// Set Cache
func (c Caches) Set(key string, value interface{}, expire ...time.Duration) {
	if _, ok := c[key]; !ok {
		c[key] = make(Cache)
	}
	var _expire time.Duration
	if len(expire) == 0 {
		_expire = DefaultMemCacheExpiration
	} else {
		_expire = expire[0]
	}
	c[key]["expire"] = time.Now().Unix() + int64(_expire)
	c[key]["value"] = value
}

// Get Cache
func (c Caches) Get(key string) interface{} {
	if cache, ok := c[key]; !ok {
		return nil
	} else {
		if expire, ok := cache["expire"]; !ok || expire.(int64) < time.Now().Unix() {
			delete(c, key)
			return nil
		}
		if value, ok := cache["value"]; ok {
			return value
		}
		delete(c, key)
		return nil
	}
}
