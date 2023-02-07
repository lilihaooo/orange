package redis

import (
	"github.com/garyburd/redigo/redis"
	"github.com/lilihaooo/orange/settings"
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
			res, err := redis.Dial(settings.Conf.RedisConfig.Protocol, address, redis.DialPassword(settings.Conf.RedisConfig.Password))
			res.Do("select db" + settings.Conf.RedisConfig.DB)
			return res, err
		},
	}
}
