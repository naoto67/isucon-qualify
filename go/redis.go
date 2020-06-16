package main

import (
	"time"

	"github.com/gomodule/redigo/redis"
)

var redisAddress = "localhost:6379"
var redisPool *redis.Pool

func newRedis() *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		MaxActive:   0,
		IdleTimeout: 240 * time.Second,
		Dial:        func() (redis.Conn, error) { return redis.Dial("tcp", redisAddress) },
	}
}
