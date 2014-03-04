package scache

import (
	"github.com/karlseguin/gspec"
	"github.com/karlseguin/nd"
	"testing"
	"time"
)

func TestGetReturnsANilValueOnMiss(t *testing.T) {
	spec := gspec.New(t)
	cache := New(Configure())
	spec.Expect(cache.Get("not valid")).ToBeNil()
}

func TestGetReturnsANilOnExpiredValue(t *testing.T) {
	spec := gspec.New(t)
	cache := New(Configure().TTL(time.Second * 10))
	cache.Set("expired", 123)
	nd.ForceNow(time.Now().Add(time.Second * 11))
	spec.Expect(cache.Get("expired")).ToBeNil()
}

func TestGetReturnsAValidValue(t *testing.T) {
	spec := gspec.New(t)
	cache := New(Configure().TTL(time.Second * 10))
	cache.Set("valid", 123)
	nd.ForceNow(time.Now().Add(time.Second * 9))
	spec.Expect(cache.Get("valid").(int)).ToEqual(123)
}

func TestSetOverwritesAnExistingValue(t *testing.T) {
	spec := gspec.New(t)
	cache := New(Configure())
	cache.Set("valid", 1)
	cache.Set("valid", 2)
	spec.Expect(cache.Get("valid").(int)).ToEqual(2)
}
