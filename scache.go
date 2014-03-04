// An LRU cache aimed at holding a small set of values
package scache

import (
	"github.com/karlseguin/nd"
	"math/rand"
	"sort"
	"sync"
	"sync/atomic"
	"time"
)

type Cache struct {
	*Configuration
	sync.RWMutex
	times   Times
	scratch []string
	lookup  map[string]*Item
}

type Item struct {
	key      string
	value    interface{}
	expires  int64
	accessed int64
}

type Times []int64

func (t Times) Len() int {
	return len(t)
}

func (t Times) Less(i, j int) bool {
	return t[i] < t[j]
}

func (t Times) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

func New(config *Configuration) *Cache {
	c := &Cache{
		Configuration: config,
		lookup:        make(map[string]*Item),
		times:         make(Times, config.workSize),
		scratch:       make([]string, config.workSize),
	}
	go c.gc()
	return c
}

func (c *Cache) Get(key string) interface{} {
	c.RLock()
	item, ok := c.lookup[key]
	c.RUnlock()
	if ok == false {
		return nil
	}
	now := nd.Now().Unix()
	if atomic.LoadInt64(&item.expires) < now {
		return nil
	}
	atomic.StoreInt64(&item.accessed, now)
	return item.value
}

func (c *Cache) Set(key string, value interface{}, ttl time.Duration) {
	now := nd.Now()
	item := &Item{
		key:      key,
		value:    value,
		accessed: now.Unix(),
		expires:  now.Add(ttl).Unix(),
	}
	c.Lock()
	c.lookup[key] = item
	c.Unlock()
}

func (c *Cache) Fetch(key string, ttl time.Duration, fetch func() (interface{}, error)) (interface{}, error) {
	item := c.Get(key)
	if item != nil {
		return item, nil
	}
	value, err := fetch()
	if err == nil {
		c.Set(key, value, ttl)
	}
	return value, err
}

func (c *Cache) gc() {
	for {
		time.Sleep(c.pruneFrequency)
		for {
			c.RLock()
			l := len(c.lookup)
			c.RUnlock()
			if l > c.maxItems {
				c.prune()
			} else {
				break
			}
		}
	}
}

func (c *Cache) prune() {
	c.RLock()
	l := int32(len(c.lookup))
	found := 0
	for _, item := range c.lookup {
		if rand.Int31n(l) > int32(c.workSize) {
			continue
		}
		c.times[found] = atomic.LoadInt64(&item.accessed)
		found++
		if found == c.workSize {
			break
		}
	}
	c.RUnlock()

	sort.Sort(c.times[:found])
	target := c.times[0]

	c.RLock()
	found = 0
	for key, item := range c.lookup {
		if atomic.LoadInt64(&item.accessed) > target {
			continue
		}
		c.scratch[found] = key
		found++
		if found == c.workSize {
			break
		}
	}
	c.RUnlock()

	c.Lock()
	for _, key := range c.scratch[:found] {
		delete(c.lookup, key)
	}
	c.Unlock()
}
