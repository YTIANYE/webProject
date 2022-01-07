package model

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
)

// 存储图片 id 到 redis 数据库
func SaveImgCode(code, uuid string) error {
	// 1. 连接数据库
	conn, err := redis.Dial("tcp", "192.168.17.129:6379")
	if err != nil {
		fmt.Println("redis Dial err:", err)
		return err
	}
	defer conn.Close()

	// 2. 写入数据库  --- 有效时间 5分钟 ，
	_, err = conn.Do("setex", uuid, 60*5, code) //setex 需要设置有效时长

	return err // 不需要回复助手
}
