package hashicorp

import "github.com/hashicorp/golang-lru"

type LRUCache struct {
    cache *lru.Cache
}

func (c *LRUCache) Get(key interface{}) interface{} {
    v, _ := c.cache.Get(key)
    return v
}

func (c *LRUCache) Put(key, value interface{}) {
    c.cache.Add(key, value)
}
