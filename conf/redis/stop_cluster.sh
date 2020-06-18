#! /bin/sh

for p in 7000 7001 7002 7003 7004 7005
do
  redis-cli -p p shutdown
done
