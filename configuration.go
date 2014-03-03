package scache

import (
	"time"
)

// Configuration options for the cache
type Configuration struct {
	workSize       int
	maxItems       int
	ttl            time.Duration
	pruneFrequency time.Duration
}

func Configure() *Configuration {
	return &Configuration{
		workSize:       50,
		maxItems:       1000,
		ttl:            time.Hour,
		pruneFrequency: time.Minute * 5,
	}
}

// The maximum amount of items which the cache should hold. Since purging is
// scheduled, the actual number can grow much larger. If this is a problem,
// this cache isn't the right solution for you.
//
// [1000]
func (c *Configuration) MaxItems(count int) *Configuration {
	c.maxItems = count
	c.workSize = count / 20
	return c
}

// How long to cache items. On an expired GET, the item will automatically
// be purged
//
// [time.Hour]
func (c *Configuration) TTL(ttl time.Duration) *Configuration {
	c.ttl = ttl
	return c
}

// The frequency to schedule a pruning
//
// [time.Minute * 5]
func (c *Configuration) PruneFrequency(frequency time.Duration) *Configuration {
	c.pruneFrequency = frequency
	return c
}
