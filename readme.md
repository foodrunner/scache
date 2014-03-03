# Overview
`scache` is an LRU cached for holding a small number of values. Rather than using
a traditional hash + list, `scache` only uses a hash. Eviction happens by sampling
the hash for a reasonably old timestamp and purging all older items.

Furthermore, since the number and size of objects is expected to be small,
`scache` purging routine runs on a schedule (rather than on-demand). It's possible
for the number of items to exceed the specified maximum.

For a more powerful and idiomatic LRU cache, check out
[ccache](https://github.com/karlseguin/ccache)

## Usage

  import (
    "github.com/foodrunner/scache"
  )

  cache := scache.New(scache.Configure().TTL(time.Hour).MaxItems(2000))

  cache.Get("goku") //nil
  cache.Set("goku", 9000)
  cache.Get("goku") // 9000
