package redis

import (
	"github.com/garyburd/redigo/redis"
	"orange/settings"
)

var Pool *redis.Pool
var KeyPrefix = "orange_"

func InitRedisConn() {
	Pool = &redis.Pool{
		MaxIdle:     16,
		MaxActive:   0,
		IdleTimeout: 300,
		Dial: func() (redis.Conn, error) {
			address := settings.Conf.RedisConfig.Host + ":" + settings.Conf.RedisConfig.Port
			return redis.Dial("tcp", address, redis.DialPassword(settings.Conf.RedisConfig.Password))

		},
	}
}
