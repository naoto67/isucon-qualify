package main

import (
	"time"

	"github.com/chasex/redis-go-cluster"
)

var redisAddress = "localhost:6379"

// var redisPool *redis.Pool
var redisCluster redis.Cluster

func newRedis() (redis.Cluster, error) {
	return redis.NewCluster(
		&redis.Options{
			StartNodes:   []string{"127.0.0.1:7000", "127.0.0.1:7001", "127.0.0.1:7002", "127.0.0.1:7003", "127.0.0.1:7004", "127.0.0.1:7005"},
			ConnTimeout:  50 * time.Millisecond,
			ReadTimeout:  50 * time.Millisecond,
			WriteTimeout: 50 * time.Millisecond,
			KeepAlive:    16,
			AliveTime:    60 * time.Second,
		})
}

func FLUSH_ALL() {
	redisCluster.Do("FLUSHALL")
}
