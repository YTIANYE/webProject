package model

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"strconv"
	"time"
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

// 创建订单

type UserData struct {
	Id int
}

func CreateOrder(houseId, startDate, endDate, userName string) (int, error) {
	var order OrderHouse

	// 传入的信息
	hid, _ := strconv.Atoi(houseId)
	order.HouseId = uint(hid)
	bDate, _ := time.Parse("2006-01-02", startDate)
	order.Begin_date = bDate
	eDate, _ := time.Parse("2006-01-02", endDate)
	order.End_date = eDate

	// 用户数据
	var userData UserData
	if err := GlobalConn.Raw("select id from user where name = ?", userName).Scan(&userData).Error; err != nil {
		fmt.Println("获取用户数据错误", err)
		return 0, err
	}
	order.UserId = uint(userData.Id)

	// 状态
	order.Status = "WAIT_ACCEPT"

	// 预定总天数
	dur := eDate.Sub(bDate)
	order.Days = int(dur.Hours()) / 24

	// 房屋单价和总价
	var house House
	GlobalConn.Where("id = ?", hid).Find(&house).Select("price")
	order.House_price = house.Price
	order.Amount = house.Price * order.Days

	// 创建订单
	if err:= GlobalConn.Create(&order).Error;err != nil{
		fmt.Println("创建订单失败：", err)
		return 0, err
	}
	return int(order.ID), nil
}
