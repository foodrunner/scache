# Overview
`scache` is an LRU cached for holding a small number of values. Rather than using a traditional hash + list, `scache` only uses a hash. Eviction happens by sampling the hash for a reasonably old timestamp and purging all older items.

For a more powerful and idiomatic LRU cache, check out [ccache](https://github.com/karlseguin/ccache)

## Use Case
`scache` has a fairly narrow, yet frequently needed, use case. It's meant to cache a relatively
small and static number of resources, while providing some basic LRU/cleanup facilities. For example, we use it to cache internal users by their authentication token. In the course of a day, we might only see a few hundred unique tokens=>users. Simply sticking them in a hashtable will eventually leak.

It's important to realize that pruning the cache happens at a confiruable interval. There's no upper limit to how many items the cache will actually store. The only guarantee is that, after a prune, the size will be less than the configurable `MaxItems`. In other words, if you expect huge spikes of entries and are concerned about memory, don't use this cache.

## Usage

    import (
      "github.com/foodrunner/scache"
    )

    cache := scache.New(scache.Configure().MaxItems(2000))

    cache.Get("goku") //nil
    cache.Set("goku", 9000, time.Hour * 2)
    cache.Get("goku") // 9000

    cache.Fetch("Leto", time.Minute * 30, func() (interface{}, error) {
      //db.load...
      // return res, err
    })


## Configuration Options
- `MaxItems` maximum items to keep in the cache
- `PruneFrequency` pruner's frequency
