#! /bin/sh

redis-server /etc/redis/cluster7000.conf &&
redis-server /etc/redis/cluster7001.conf &&
redis-server /etc/redis/cluster7002.conf &&
redis-server /etc/redis/cluster7003.conf &&
redis-server /etc/redis/cluster7004.conf &&
redis-server /etc/redis/cluster7005.conf

redis-cli --cluster create 127.0.0.1:7000 127.0.0.1:7001 127.0.0.1:7002 127.0.0.1:7003 127.0.0.1:7004 127.0.0.1:7005 --cluster-replicas 1
