package main

import (
	"context"

	"github.com/go-redis/redis/v8"
)

var redisCtx = context.Background()

// var redisPool *redis.Pool
var redisCluster *redis.ClusterClient

func newRedis() (*redis.ClusterClient, error) {
	rdb := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: []string{":7006", ":7001", ":7002", ":7003", ":7004", ":7005"},
	})
	return rdb, nil
}

func FLUSH_ALL() {
	redisCluster.Do(redisCtx, "FLUSHALL")
}
