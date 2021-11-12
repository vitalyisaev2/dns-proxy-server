package lru

import (
    "github.com/hashicorp/golang-lru"
    "github.com/mageddo/dns-proxy-server/cache"
)

type LRUCache struct {
    Cache *lru.Cache
}

func (c *LRUCache) ContainsKey(key interface{}) bool {
    return c.Cache.Contains(key)
}

func (c *LRUCache) Get(key interface{}) interface{} {
    v, _ := c.Cache.Get(key)
    return v
}

//
// Put value in cache, it doesn't have guarantee of concurrency treat
//
func (c *LRUCache) Put(key, value interface{}) {
    c.Cache.Add(key, value)
}

//
// Check if value is already associated, if yes just return it, if not put the passed value and return nil
// This method must be thread safe (atomic)
//
func (c *LRUCache) PutIfAbsent(key, value interface{}) interface{} {
    if ok, _ := c.Cache.ContainsOrAdd(key, value); ok {
        return c.Get(key)
    }
    return nil
}

func (c *LRUCache) Clear() {
    c.Cache.Purge()
}

func (c *LRUCache) Remove(key interface{}) {
    c.Cache.Remove(key)
}

func (c *LRUCache) KeySet() []interface{} {
    return c.Cache.Keys()
}

func (c *LRUCache) Size() int {
    return c.Cache.Len()
}

//
// Creates a LRU cache
// size is the maximum size of the cache, -1 if it is unlimited
//
func New(size int) cache.Cache {
    c, err := lru.New(size)
    if err != nil {
        return nil
    }
    return &LRUCache{c}
}
