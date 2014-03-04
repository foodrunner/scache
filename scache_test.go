package scache

import (
	"errors"
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
	cache := New(Configure())
	cache.Set("expired", 123, time.Second * 10)
	nd.ForceNow(time.Now().Add(time.Second * 11))
	spec.Expect(cache.Get("expired")).ToBeNil()
}

func TestGetReturnsAValidValue(t *testing.T) {
	spec := gspec.New(t)
	cache := New(Configure())
	cache.Set("valid", 123, time.Second * 10)
	nd.ForceNow(time.Now().Add(time.Second * 9))
	spec.Expect(cache.Get("valid").(int)).ToEqual(123)
}

func TestClearErasesTheCache(t *testing.T) {
	spec := gspec.New(t)
	cache := New(Configure())
	cache.Set("valid", 123, time.Second * 10)
	cache.Set("valid2", 55, time.Second * 10)
	spec.Expect(len(cache.lookup)).ToEqual(2)
	cache.Clear()
	spec.Expect(len(cache.lookup)).ToEqual(0)
}

func TestSetOverwritesAnExistingValue(t *testing.T) {
	spec := gspec.New(t)
	cache := New(Configure())
	cache.Set("valid", 1, time.Hour)
	cache.Set("valid", 2, time.Hour)
	spec.Expect(cache.Get("valid").(int)).ToEqual(2)
}

func TestFetchReturnsAValidValue(t *testing.T) {
	spec := gspec.New(t)
	cache := New(Configure())
	cache.Set("valid", 123, time.Hour)
	value, err := cache.Fetch("valid", time.Hour, nil)
	spec.Expect(err).ToBeNil()
	spec.Expect(value.(int)).ToEqual(123)
}

func TestFetchLoadsOnMiss(t *testing.T) {
	spec := gspec.New(t)
	cache := New(Configure())
	value, err := cache.Fetch("miss", time.Hour, func() (interface{}, error) {
		return 14495, nil
	})
	spec.Expect(err).ToBeNil()
	spec.Expect(value.(int)).ToEqual(14495)
	spec.Expect(cache.Get("miss").(int)).ToEqual(14495)
}

func TestFetchPassesError(t *testing.T) {
	spec := gspec.New(t)
	cache := New(Configure())
	value, err := cache.Fetch("miss", time.Minute, func() (interface{}, error) {
		return nil, errors.New("fetch fail")
	})
	spec.Expect(err.Error()).ToEqual("fetch fail")
	spec.Expect(value).ToBeNil()
}
