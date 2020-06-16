package main

import (
	"github.com/bradfitz/gomemcache/memcache"
)

const memcachedAddress = "localhost:11211"

type CacheClient struct {
	cli *memcache.Client
}

var cacheClient CacheClient

func NewCacheClient() {
	cli := memcache.New(memcachedAddress)
	cli.MaxIdleConns = 10
	cacheClient = CacheClient{cli: cli}
}

func (c CacheClient) Set(key string, value []byte) error {
	return c.cli.Set(&memcache.Item{Key: key, Value: value})
}

func (c CacheClient) Get(key string) ([]byte, error) {
	item, err := c.cli.Get(key)
	if err != nil {
		return nil, err
	}
	return item.Value, nil
}

func (c CacheClient) FLUSH() {
	c.cli.FlushAll()
}
