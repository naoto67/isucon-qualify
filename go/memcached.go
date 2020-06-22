package main

import (
	"github.com/bradfitz/gomemcache/memcache"
)

const memcached1Address = "localhost:11211"
const memcached2Address = "localhost:11212"

type CacheClient struct {
	cli1 *memcache.Client
	cli2 *memcache.Client
}

var cacheClient CacheClient

func NewCacheClient() {
	cli1 := memcache.New(memcached1Address)
	cli2 := memcache.New(memcached2Address)
	cli1.MaxIdleConns = 10
	cli2.MaxIdleConns = 10
	cacheClient = CacheClient{cli1: cli1, cli2: cli2}
}

func (c CacheClient) Set(key string, value []byte) error {
	if len(key)%2 == 0 {
		return c.cli1.Set(&memcache.Item{Key: key, Value: value})
	} else {
		return c.cli2.Set(&memcache.Item{Key: key, Value: value})
	}
}

func (c CacheClient) Get(key string) ([]byte, error) {
	var item *memcache.Item
	var err error
	if len(key)%2 == 0 {
		item, err = c.cli1.Get(key)
	} else {
		item, err = c.cli2.Get(key)
	}
	if err != nil {
		return nil, err
	}
	return item.Value, nil
}

func (c CacheClient) Delete(key string) error {
	if len(key)%2 == 0 {
		return c.cli1.Delete(key)
	} else {
		return c.cli2.Delete(key)
	}
}

func (c CacheClient) FLUSH() {
	c.cli1.FlushAll()
	c.cli2.FlushAll()
}
