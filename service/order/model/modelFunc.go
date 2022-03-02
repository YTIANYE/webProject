package model

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"service/order/proto/order"
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
	if err := GlobalConn.Create(&order).Error; err != nil {
		fmt.Println("创建订单失败：", err)
		return 0, err
	}
	return int(order.ID), nil
}

// 查看用户订单

func GetOrderInfo(role, userName string) ([]*order.OrdersData, error) {
	// 返回结果
	var orderResp []*order.OrdersData

	// 用户的所有订单
	var orders []OrderHouse

	var userData UserData
	//用原生查询的时候,查询的字段必须跟数据库中的字段保持一直
	GlobalConn.Raw("select id from user where name = ?",userName).Scan(&userData)

	// 用户 分为 顾客 和 房主
	if role == "custom"{//顾客
		if err := GlobalConn.Where("user_id = ?", userData.Id).Find(&orders).Error; err != nil{
			fmt.Println("获取顾客的订单信息失败", err)
			return nil, err
		}
	}else {// 房主
		var houses []House
		GlobalConn.Where("user_id = ?", userData.Id).Find(&houses)
		for  _,v := range houses{
			var tempOrders []OrderHouse
			GlobalConn.Model(&v).Related(tempOrders)
			orders = append(orders, tempOrders...)
		}
	}

	// 添加返回信息
	for _, v := range orders{
		var orderTemp order.OrdersData
		orderTemp.OrderId = int32(v.ID)
		orderTemp.EndDate = v.End_date.Format("2006-01-02")
		orderTemp.StartDate = v.Begin_date.Format("2006-01-02")
		orderTemp.Ctime = v.CreatedAt.Format("2006-01-02")
		orderTemp.Amount = int32(v.Amount)
		orderTemp.Comment = v.Comment
		orderTemp.Days = int32(v.Days)
		orderTemp.Status = v.Status

		//关联house表
		var house House
		GlobalConn.Model(&v).Related(&house).Select("index_image_url","title")
		orderTemp.ImgUrl = "http://192.168.17.129:8888/"+house.Index_image_url
		orderTemp.Title = house.Title

		orderResp = append(orderResp, &orderTemp)
	}
	return orderResp,nil
}
