package model

import (
	"github.com/gomodule/redigo/redis"
)

// 创建全局连接池 redis 句柄

var RedisPool redis.Pool

// 创建函数 初始化连接池

func InitRedis() {
	// 连接 Redis 连接池
	RedisPool = redis.Pool{
		MaxIdle:         20,
		MaxActive:       50,
		MaxConnLifetime: 60 * 5,
		IdleTimeout:     60,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", "192.168.17.129:6379")
		},
	}
}
