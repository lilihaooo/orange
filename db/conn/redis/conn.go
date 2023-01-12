package redis

import (
	"github.com/garyburd/redigo/redis"
)

var Pool *redis.Pool
var KeyPrefix = "orange"

func InitRedisConn() {
	Pool = &redis.Pool{
		MaxIdle:     16,
		MaxActive:   0,
		IdleTimeout: 300,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", "127.0.0.1:6379", redis.DialPassword("123456"))
		},
	}
}
