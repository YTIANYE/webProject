package model

import (
	"fmt"
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

// 校验图片验证码
func CheckImgCode(uuid, imgCode string) bool {
	// 连接redis  ---  从连接池中获取链接
	/*	conn, err := redis.Dial("tcp", "192.168.17.129:6379")
		if err != nil {
			fmt.Println("redis.Dial err:", err)
		}*/
	conn := RedisPool.Get()
	defer conn.Close()

	// 查询redis中的数据
	code, err := redis.String(conn.Do("get", uuid))
	if err != nil {
		fmt.Println("查询错误 err:", err)
		return false
	}

	// 返回校验结果
	return code == imgCode
}

// 存储短信验证码
func SaveSmsCode(phone, code string) error {
	// 从 redis 连接池获取一条连接
	conn := RedisPool.Get()
	defer conn.Close() //相当于释放链接，不是关闭链接

	// 存储短信验证码 到redis中
	_, err := conn.Do("setex", phone+"_code", 60*3, code)
	return err
}
